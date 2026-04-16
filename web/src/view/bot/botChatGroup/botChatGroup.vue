<template>
  <div>
    <!-- 搜索 -->
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="创建日期">
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            type="datetimerange"
          />
        </el-form-item>

        <el-form-item label="机器人">
          <el-select v-model="searchInfo.botID" clearable>
            <el-option
              v-for="item in botOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="群名称">
          <el-input v-model="searchInfo.chatGroupName" clearable />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格 -->
    <div class="gva-table-box">
      <div style="margin-bottom: 10px;">
        <el-button
          type="success"
          :disabled="!multipleSelection.length"
          @click="openBindClassify"
        >
          加入群聊分组
        </el-button>

        <el-button
          type="danger"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >
          删除选中群聊
        </el-button>
      </div>

      <el-table :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column label="ID" prop="ID" />
        <el-table-column label="机器人" prop="botName" />
        <el-table-column label="群组名称" prop="chatGroupName" />

        <el-table-column label="存储群消息">
          <template #default="{ row }">
            <el-switch
              :model-value="row.syncMessage === 1"
              @change="val => onSyncChange(val, row)"
            />
          </template>
        </el-table-column>

        <el-table-column label="禁用转发">
          <template #default="{ row }">
            <el-switch
              :model-value="row.banForward === 1"
              @change="val => onBanForwardChange(val, row)"
            />
          </template>
        </el-table-column>

        <el-table-column label="禁用文本长度">
          <template #default="{ row }">
            <el-input-number
              v-model="row.maxWords"
              @change="val => onMaxWordsChange(val, row)"
            />
          </template>
        </el-table-column>
      </el-table>

      <!-- ✅ 分页（已增强） -->
      <el-pagination
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        :page-size="pageSize"
        :current-page="page"
        :page-sizes="[20, 50, 100, 200]"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>

    <!-- 弹窗 -->
    <el-dialog v-model="bindDialogVisible" title="加入群聊分组" width="500px">
      <div style="margin-bottom: 20px;">
        <el-input
          v-model="newClassifyTitle"
          placeholder="输入新分组名称"
          style="width: 70%; margin-right: 10px;"
        />
        <el-button type="primary" @click="submitCreateClassify">
          创建分组
        </el-button>
      </div>

      <el-divider>选择已有分组</el-divider>

      <el-radio-group v-model="selectedClassifyID">
        <el-radio
          v-for="item in classifyOptions"
          :key="item.ID"
          :label="item.ID"
          style="display: block; margin-bottom: 10px;"
        >
          {{ item.title }}
        </el-radio>
      </el-radio-group>

      <template #footer>
        <el-button @click="bindDialogVisible=false">取消</el-button>
        <el-button type="primary" @click="submitBindClassify">
          确定
        </el-button>
      </template>
    </el-dialog>

  </div>
</template>
<script setup>
import {
  getBotChatGroupList,
  updateBotChatGroup,
  deleteBotChatGroupByIds,
  saveBotChatGroupClassify,
  chooseChatGroupClassify
} from '@/api/bot/botChatGroup'

import { getBotChoice } from '@/api/bot/bot'
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

/* 多选 */
const multipleSelection = ref([])
const handleSelectionChange = val => {
  multipleSelection.value = val
}

/* 删除 */
const onDelete = async () => {
  const ids = multipleSelection.value.map(i => i.ID)
  const res = await deleteBotChatGroupByIds({ ids })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    getTableData()
  }
}

/* 机器人 */
const botOptions = ref([])
const loadBots = async () => {
  const res = await getBotChoice()
  if (res.code === 0) {
    botOptions.value = res.data.map(i => ({
      label: i.name,
      value: i.botID
    }))
  }
}
loadBots()

/* 表格 */
const tableData = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const searchInfo = ref({})

const getTableData = async () => {
  const res = await getBotChatGroupList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}
onMounted(getTableData)

/* 行更新 */
const updateRow = async (row, patch) => {
  await updateBotChatGroup({ ...row, ...patch })
  Object.assign(row, patch)
}

const onSyncChange = (val, row) => updateRow(row, { syncMessage: val ? 1 : 2 })
const onBanForwardChange = (val, row) => updateRow(row, { banForward: val ? 1 : 2 })
const onMaxWordsChange = (val, row) => updateRow(row, { maxWords: val })

/* 分组 */
const classifyOptions = ref([])
const newClassifyTitle = ref('')
const bindDialogVisible = ref(false)
const selectedClassifyID = ref(null)

const loadClassify = async () => {
  const res = await chooseChatGroupClassify()
  if (res.code === 0) {
    classifyOptions.value = res.data.list || []
  }
}

const openBindClassify = async () => {
  bindDialogVisible.value = true
  await loadClassify()
}

const submitCreateClassify = async () => {
  if (!newClassifyTitle.value) {
    ElMessage.warning('请输入分组名称')
    return
  }
  const res = await saveBotChatGroupClassify({ title: newClassifyTitle.value })
  if (res.code === 0) {
    ElMessage.success('创建成功')
    newClassifyTitle.value = ''
    await loadClassify()
  }
}

/* ✅ 核心修改：这里 */
const submitBindClassify = async () => {
  if (!selectedClassifyID.value) {
    ElMessage.warning('请选择分组')
    return
  }

  // ✅ 拼接 botID_groupID
  const ids = multipleSelection.value.map(i => {
    const botID = i.botID || 0
    const groupID = i.chatGroupID
    return `${botID}_${groupID}`
  })

  const res = await saveBotChatGroupClassify({
    ID: selectedClassifyID.value,
    chatGroups: ids.join(','), // ✅ 改这里
    refresh: false
  })

  if (res.code === 0) {
    ElMessage.success('操作成功')
    bindDialogVisible.value = false
  }
}

/* 分页 */
const handleCurrentChange = p => {
  page.value = p
  getTableData()
}

const handleSizeChange = size => {
  pageSize.value = size
  page.value = 1
  getTableData()
}

/* 搜索 */
const onSubmit = () => {
  page.value = 1
  getTableData()
}
const onReset = () => {
  searchInfo.value = {}
  page.value = 1
  getTableData()
}
</script>