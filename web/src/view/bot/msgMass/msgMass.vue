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
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">

        <el-form-item label="机器人">
          <el-select
            v-model="searchInfo.botID"
            placeholder="请选择机器人"
            clearable
            style="width: 200px"
            @change="onSearchBotChange"
          >
            <el-option
              v-for="bot in botList"
              :key="bot.botID"
              :label="bot.name"
              :value="bot.botID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="群聊">
          <el-select
            v-model="searchInfo.chatGroupID"
            placeholder="请选择群聊"
            clearable
            style="width: 260px"
            :disabled="!searchInfo.botID"
          >
            <el-option
              v-for="g in searchChatGroups"
              :key="g.chatGroupID"
              :label="g.chatGroupName"
              :value="g.chatGroupID"
            />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>

      </el-form>
    </div>

    <!-- ================= 表格 ================= -->
    <div class="gva-table-box">
      <el-table
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="机器人" prop="botName" width="180" />
        <el-table-column label="群聊" prop="chatGroupName" width="180" />
        <el-table-column label="成员" prop="members" min-width="300" />
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
      <template #header>{{ type === 'create' ? '新增' : '编辑' }}</template>

      <el-form label-position="top">
        <el-form-item label="机器人">
          <el-select v-model="formData.botID" @change="onBotChange" style="width:100%">
            <el-option
              v-for="bot in botList"
              :key="bot.botID"
              :label="bot.name"
              :value="bot.botID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="群聊">
          <el-select
            v-model="formData.chatGroupID"
            :disabled="!formData.botID"
            style="width:100%"
          >
            <el-option
              v-for="g in currentChatGroups"
              :key="g.chatGroupID"
              :label="g.chatGroupName"
              :value="g.chatGroupID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="成员">
          <el-input v-model="formData.members" type="textarea" :rows="5" />
        </el-form-item>
      </el-form>

      <div style="text-align:right">
        <el-button @click="closeDialog">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </div>
    </el-drawer>

    <!-- ================= 编辑发送内容 ================= -->
    <el-dialog v-model="sendDialogVisible" title="编辑发送内容" width="70%">
      <RichEdit v-model="sendContent" />

      <el-form style="margin-top:10px">
        <el-form-item label="是否艾特">
          <el-switch v-model="atUsers" active-text="是" inactive-text="否" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="sendDialogVisible=false">取消</el-button>
        <el-button type="primary" @click="sendDialogVisible=false">确认</el-button>
      </template>
    </el-dialog>

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

const tableData = ref([])
const page = ref(1)
const pageSize = ref(30)
const total = ref(0)
const multipleSelection = ref([])

const sendContent = ref('')
const atUsers = ref(false)
const sendDialogVisible = ref(false)

const dialogVisible = ref(false)
const detailVisible = ref(false)
const type = ref('create')

const formData = ref({
  ID: null,
  botID: '',
  chatGroupID: '',
  members: ''
})

const detailForm = ref({})

const botList = ref([])
let botLoaded = false

const loadBotListLazy = async () => {
  if (botLoaded) return
  const res = await getBotChoiceWithChatGroup()
  if (res.code === 0) {
    botList.value = res.data
    botLoaded = true
  }
}

onMounted(() => {
  loadBotListLazy()
  getTableData()
})

/* ================= 搜索 ================= */

const searchInfo = ref({
  botID: '',
  chatGroupID: ''
})

const searchChatGroups = computed(() => {
  if (!searchInfo.value.botID) return []
  const bot = botList.value.find(b => b.botID === searchInfo.value.botID)
  return bot?.botChatGroups || []
})

const onSearchBotChange = () => {
  searchInfo.value.chatGroupID = ''
}

const getTableData = async () => {
  const payload = {
    page: page.value,
    pageSize: pageSize.value
  }

  if (searchInfo.value.botID) payload.botID = searchInfo.value.botID
  if (searchInfo.value.chatGroupID) payload.chatGroupID = searchInfo.value.chatGroupID

  const res = await getBotMsgMassList(payload)
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = { botID: '', chatGroupID: '' }
  page.value = 1
  getTableData()
}

/* ================= 表格 & 编辑 ================= */

const handleSelectionChange = val => {
  multipleSelection.value = val
}

const onBotChange = () => {
  formData.value.chatGroupID = ''
}

const save = async () => {
  if (!formData.value.botID || !formData.value.chatGroupID) {
    ElMessage.warning('请选择机器人和群聊')
    return
  }
  const api = type.value === 'create' ? createBotMsgMass : updateBotMsgMass
  const res = await api(formData.value)
  if (res.code === 0) {
    ElMessage.success('保存成功')
    dialogVisible.value = false
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

const openSendDialog = () => {
  sendDialogVisible.value = true
}

const sendBatch = async () => {
  const ids = multipleSelection.value.map(i => i.ID)
  const res = await sendBotMsgMass({
    msg: sendContent.value,
    ids,
    atUsers: atUsers.value
  })
  if (res.code === 0) ElMessage.success('群发成功')
}

const getDetails = async row => {
  const res = await findBotMsgMass({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    detailVisible.value = true
  }
}
</script>
