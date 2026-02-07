<template>
  <div>

    <!-- ================= 搜索 ================= -->
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="创建日期">
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            type="datetimerange"
            range-separator="至"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

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

        <el-table-column label="日期">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>

        <el-table-column label="机器人" prop="botName" />
        <el-table-column label="群聊" prop="chatGroupName" />
        <el-table-column label="成员" prop="members" />

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

        <el-form-item label="机器人">
          <el-select v-model="formData.botID">
            <el-option
              v-for="bot in botList"
              :key="bot.botID"
              :label="bot.name"
              :value="bot.botID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="群聊">
          <el-select v-model="formData.chatGroupID">
            <el-option
              v-for="g in currentChatGroups"
              :key="g.chatGroupID"
              :label="g.chatGroupName"
              :value="g.chatGroupID"
            />
          </el-select>
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

    <!-- ================= 详情 ================= -->
    <el-drawer v-model="detailVisible" title="详情">
      <el-descriptions border column="1">
        <el-descriptions-item label="机器人">{{ detailForm.botName }}</el-descriptions-item>
        <el-descriptions-item label="群聊">{{ detailForm.chatGroupName }}</el-descriptions-item>
        <el-descriptions-item label="成员">{{ detailForm.members }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>

  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'

import {
  createBotMsgMass,
  updateBotMsgMass,
  deleteBotMsgMass,
  deleteBotMsgMassByIds,
  findBotMsgMass,
  getBotMsgMassList,
  sendBotMsgMass   // ✅ 群发接口
} from '@/api/bot/msgMass'

import { getBotChoiceWithChatGroup } from '@/api/bot/bot'
import RichEdit from '@/components/richtext/rich-edit.vue'

// ================= 表格 =================
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchInfo = ref({})
const multipleSelection = ref([])

// ================= 群发 =================
const sendContent = ref('')
const sendDialogVisible = ref(false)

// ================= 表单 =================
const dialogVisible = ref(false)
const detailVisible = ref(false)
const type = ref('create')

const formData = ref({
  botID: undefined,
  chatGroupID: undefined,
  members: ''
})

const detailForm = ref({})

// ================= 机器人 =================
const botList = ref([])

const currentChatGroups = computed(() => {
  const bot = botList.value.find(b => b.botID === formData.value.botID)
  return bot?.botChatGroups || []
})

const loadBotList = async () => {
  const res = await getBotChoiceWithChatGroup()
  if (res.code === 0) botList.value = res.data
}
loadBotList()

// ================= 表格数据 =================
const getTableData = async () => {
  const res = await getBotMsgMassList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}
getTableData()

const handleSelectionChange = val => {
  multipleSelection.value = val
}

// ================= CRUD =================
const openDialog = () => {
  type.value = 'create'
  dialogVisible.value = true
}

const closeDialog = () => {
  dialogVisible.value = false
  formData.value = { botID: undefined, chatGroupID: undefined, members: '' }
}

const save = async () => {
  const api = type.value === 'create' ? createBotMsgMass : updateBotMsgMass
  const res = await api(formData.value)
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
    formData.value = res.data
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
  const IDs = multipleSelection.value.map(i => i.ID)
  await deleteBotMsgMassByIds({ IDs })
  ElMessage.success('删除成功')
  getTableData()
}

// ================= 群发（核心） =================
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

const getDetails = async row => {
  const res = await findBotMsgMass({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    detailVisible.value = true
  }
}

const onSubmit = () => getTableData()
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}
</script>
