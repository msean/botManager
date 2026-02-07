<template>
  <div>

    <!-- ================= 群发按钮区 ================= -->
    <el-card style="margin-bottom: 15px">
      <el-button type="primary" @click="openSendDialog">
        编辑发送内容
      </el-button>

      <el-button
        type="success"
        style="margin-left: 10px"
        :disabled="!multipleSelection.length || !sendContent"
        @click="sendBatch"
      >
        发送到 Telegram
      </el-button>

      <span style="margin-left: 15px; color: #999">
        已选 {{ multipleSelection.length }} 条
      </span>
    </el-card>

    <!-- ================= 搜索 ================= -->
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        @keyup.enter="onSubmit"
      >
        <!-- 机器人+群聊级联 -->
        <el-form-item label="机器人/群聊">
          <el-cascader
            v-model="searchInfo.botChatGroup"
            :options="botCascaderOptions"
            placeholder="请选择机器人和群聊"
            clearable
            style="width: 300px"
            @focus="loadBotListLazy"
          />
        </el-form-item>

        <!-- 日期 -->
        <el-form-item label="创建日期">
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- ================= 表格 ================= -->
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" @click="openDialog">新增</el-button>
        <el-button
          type="danger"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >
          删除
        </el-button>
      </div>

      <el-table
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="机器人" prop="botName" width="120" />
        <el-table-column label="群聊" prop="chatGroupName" width="120" />
        <el-table-column label="成员" prop="members" width="200" />
        <el-table-column label="发送内容" width="150">
          <template #default="scope">
            <el-button link type="primary" @click="openMsg(scope.row.msg)">
              查看
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button link @click="getDetails(scope.row)">查看</el-button>
            <el-button link @click="updateRow(scope.row)">编辑</el-button>
            <el-button link type="danger" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        layout="total, sizes, prev, pager, next"
        :total="total"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        @current-change="getTableData"
        @size-change="getTableData"
      />
    </div>

    <!-- ================= 新增 / 编辑 ================= -->
    <el-drawer v-model="dialogVisible" size="40%">
      <template #header>
        <span>{{ type === 'create' ? '新增' : '编辑' }}</span>
      </template>

      <el-form :model="formData" label-position="top">
        <el-form-item label="机器人/群聊">
          <el-cascader
            v-model="formData.botChatGroup"
            :options="botCascaderOptions"
            placeholder="请选择机器人和群聊"
            clearable
            @focus="loadBotListLazy"
          />
        </el-form-item>

        <el-form-item label="成员">
          <el-input
            v-model="formData.members"
            type="textarea"
            :rows="5"
          />
        </el-form-item>

      </el-form>

      <div style="text-align: right">
        <el-button @click="closeDialog">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </div>
    </el-drawer>

    <!-- ================= 编辑发送内容弹窗 ================= -->
    <el-dialog
      v-model="sendDialogVisible"
      title="编辑发送内容"
      width="70%"
      destroy-on-close
    >
      <RichEdit v-model="sendContent" />

      <template #footer>
        <el-button @click="sendDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="sendDialogVisible = false">
          确认
        </el-button>
      </template>
    </el-dialog>

    <!-- ================= 查看发送内容 ================= -->
    <el-drawer
      v-model="contentVisible"
      title="发送内容"
      size="50%"
      destroy-on-close
    >
      <RichView v-model="currentContent" />
    </el-drawer>

    <!-- ================= 详情 ================= -->
    <el-drawer v-model="detailVisible" title="详情">
      <el-descriptions border column="1">
        <el-descriptions-item label="机器人">{{ detailForm.botName }}</el-descriptions-item>
        <el-descriptions-item label="群聊">{{ detailForm.chatGroupName }}</el-descriptions-item>
        <el-descriptions-item label="成员">{{ detailForm.members }}</el-descriptions-item>
        <el-descriptions-item label="发送内容">
          <RichView v-model="detailForm.msg" />
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import RichEdit from '@/components/richtext/rich-edit.vue'
import RichView from '@/components/richtext/rich-view.vue'

import {
  createBotMsgMass,
  updateBotMsgMass,
  deleteBotMsgMass,
  deleteBotMsgMassByIds,
  findBotMsgMass,
  getBotMsgMassList,
  sendBotMsgMass
} from '@/api/bot/msgMass'

import { getBotChoiceWithChatGroup } from '@/api/bot/bot'

// ================= 表格 =================
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const multipleSelection = ref([])

// ================= 群发 =================
const sendContent = ref('')
const sendDialogVisible = ref(false)

// ================= 弹窗 =================
const dialogVisible = ref(false)
const detailVisible = ref(false)
const type = ref('create')

const formData = ref({
  botChatGroup: [],
  members: ''
})

const detailForm = ref({})

// ================= 机器人/群聊级联 =================
const botList = ref([])
const botCascaderOptions = ref([])
let botLoaded = false
const loadBotListLazy = async () => {
  if (botLoaded) return
  const res = await getBotChoiceWithChatGroup()
  if (res.code === 0) {
    botList.value = res.data
    botCascaderOptions.value = botList.value.map(bot => ({
      value: bot.botID,
      label: bot.name,
      children: (bot.botChatGroups || []).map(g => ({
        value: g.chatGroupID,
        label: g.chatGroupName
      }))
    }))
    botLoaded = true
  }
}

// ================= 页面加载立即获取表格数据 =================
onMounted(() => {
  getTableData()
})

// ================= 表格数据 =================
const searchInfo = ref({
  botChatGroup: [],
  createdAtRange: []
})

const getTableData = async () => {
  const payload = {
    page: page.value,
    pageSize: pageSize.value,
  }
  if (searchInfo.value.botChatGroup?.length) {
    payload.botID = searchInfo.value.botChatGroup[0]
    payload.chatGroupID = searchInfo.value.botChatGroup[1]
  }
  if (searchInfo.value.createdAtRange?.length) {
    payload.startTime = searchInfo.value.createdAtRange[0]
    payload.endTime = searchInfo.value.createdAtRange[1]
  }
  const res = await getBotMsgMassList(payload)
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

const handleSelectionChange = val => {
  multipleSelection.value = val
}

// ================= CRUD =================
const openDialog = () => {
  type.value = 'create'
  formData.value = { botChatGroup: [], members: '' }
  dialogVisible.value = true
}

const closeDialog = () => {
  dialogVisible.value = false
}

const save = async () => {
  if (!formData.value.botChatGroup?.length) {
    ElMessage.warning('请选择机器人和群聊')
    return
  }
  const payload = {
    ...formData.value,
    botID: formData.value.botChatGroup[0],
    chatGroupID: formData.value.botChatGroup[1]
  }
  const api = type.value === 'create' ? createBotMsgMass : updateBotMsgMass
  const res = await api(payload)
  if (res.code === 0) {
    ElMessage.success('保存成功')
    closeDialog()
    getTableData()
  }
}

const updateRow = async row => {
  const res = await findBotMsgMass({ ID: row.ID })
  if (res.code === 0) {
    type.value = 'update'
    formData.value = {
      ...res.data,
      botChatGroup: [res.data.botID, res.data.chatGroupID]
    }
    dialogVisible.value = true
  }
}

const deleteRow = row => {
  ElMessageBox.confirm('确认删除？').then(async () => {
    await deleteBotMsgMass({ ID: row.ID })
    ElMessage.success('已删除')
    getTableData()
  })
}

const onDelete = async () => {
  if (!multipleSelection.value.length) return
  const IDs = multipleSelection.value.map(i => i.ID)
  await deleteBotMsgMassByIds({ IDs })
  ElMessage.success('删除成功')
  getTableData()
}

// ================= 群发 =================
const openSendDialog = () => {
  sendDialogVisible.value = true
}

const sendBatch = async () => {
  if (!sendContent.value) {
    ElMessage.warning('请先编辑发送内容')
    return
  }
  if (!multipleSelection.value.length) {
    ElMessage.warning('请先选择要发送的记录')
    return
  }
  const ids = multipleSelection.value.map(item => item.ID)
  const res = await sendBotMsgMass({
    msg: sendContent.value,
    ids,
    remark: ''
  })
  if (res.code === 0) {
    ElMessage.success('群发记录已保存')
  }
}

// ================= 查看发送内容 =================
const contentVisible = ref(false)
const currentContent = ref('')
const openMsg = (msg) => {
  currentContent.value = msg
  contentVisible.value = true
}

// ================= 查看详情 =================
const getDetails = async row => {
  const res = await findBotMsgMass({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    detailVisible.value = true
  }
}

// ================= 搜索/重置 =================
const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = { botChatGroup: [], createdAtRange: [] }
  getTableData()
}
</script>
