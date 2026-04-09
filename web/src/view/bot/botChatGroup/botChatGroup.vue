<template>
  <div>
    <!-- 搜索 -->
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="创建日期">
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>

```
    <el-form-item label="机器人">
      <el-select v-model="searchInfo.botID" clearable placeholder="请选择机器人">
        <el-option
          v-for="item in botOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        />
      </el-select>
    </el-form-item>

    <el-form-item label="群名称">
    <el-input 
      v-model="searchInfo.chatGroupName" 
      placeholder="请输入群名称"
      clearable
    />
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
    <el-button type="primary" @click="openCreateClassify">
      新建分组
    </el-button>

    <el-button
      type="success"
      :disabled="!multipleSelection.length"
      @click="openBindClassify"
    >
      加入分组
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
        <el-input-number v-model="row.maxWords" @change="val => onMaxWordsChange(val, row)" />
      </template>
    </el-table-column>
  </el-table>

  <el-pagination
    layout="total, prev, pager, next"
    :total="total"
    :page-size="pageSize"
    :current-page="page"
    @current-change="handleCurrentChange"
  />
</div>

<!-- 新建分组 -->
<el-dialog v-model="createDialogVisible" title="新建分组">
  <el-input v-model="newClassifyTitle" placeholder="请输入分组名称" />
  <template #footer>
    <el-button @click="createDialogVisible=false">取消</el-button>
    <el-button type="primary" @click="submitCreateClassify">确定</el-button>
  </template>
</el-dialog>

<!-- 选择分组 -->
<el-dialog v-model="bindDialogVisible" title="选择分组">
  <el-radio-group v-model="selectedClassifyID">
    <el-radio
      v-for="item in classifyOptions"
      :key="item.ID"
      :label="item.ID"
    >
      {{ item.title }}
    </el-radio>
  </el-radio-group>

  <template #footer>
    <el-button @click="bindDialogVisible=false">取消</el-button>
    <el-button type="primary" @click="submitBindClassify">确定</el-button>
  </template>
</el-dialog>
```

  </div>
</template>

<script setup>
import {
  getBotChatGroupList,
  updateBotChatGroup,
  deleteBotChatGroupByIds,
  saveBotChatGroupClassify,
  chooseChatGroupClassify,
  getBotChatGroupClassifyList
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
const pageSize = ref(10)
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

/* ===== 分组 ===== */
const createDialogVisible = ref(false)
const newClassifyTitle = ref('')

const openCreateClassify = () => {
  newClassifyTitle.value = ''
  createDialogVisible.value = true
}

const submitCreateClassify = async () => {
  const res = await saveBotChatGroupClassify({ title: newClassifyTitle.value })
  if (res.code === 0) {
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadClassify()
  }
}

/* 分组列表 */
const classifyOptions = ref([])
const loadClassify = async () => {
  const res = await chooseChatGroupClassify()
  if (res.code === 0) classifyOptions.value = res.data
}
loadClassify()

/* 加入分组 */
const bindDialogVisible = ref(false)
const selectedClassifyID = ref(null)

const openBindClassify = () => {
  bindDialogVisible.value = true
}

const submitBindClassify = async () => {
  const ids = multipleSelection.value.map(i => i.chatGroupID)

  const res = await getBotChatGroupClassifyList({
    classifyID: selectedClassifyID.value,
    chatGroupIDs: ids
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

/* 搜索 */
const onSubmit = () => getTableData()
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}
</script>
