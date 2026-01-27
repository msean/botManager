<template>
  <div>
    <!-- 搜索 -->
    <div class="gva-search-box">
      <el-form
        :inline="true"
        :model="searchInfo"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="创建日期">
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>

        <el-form-item label="机器人">
          <el-select
            v-model="searchInfo.botID"
            filterable
            clearable
            placeholder="请选择机器人"
            style="width: 220px"
          >
            <el-option
              v-for="item in botOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格 -->
    <div class="gva-table-box">
      <el-table :data="tableData" row-key="ID">
        <el-table-column label="机器人" prop="botName" width="180" />
        <el-table-column label="群组名称" prop="chatGroupName" />

        <!-- 同步消息 -->
        <el-table-column label="存储群消息" width="140">
          <template #default="{ row }">
            <el-switch
              :model-value="row.syncMessage === 1"
              @change="val => onSyncChange(val, row)"
            />
          </template>
        </el-table-column>

        <!-- 禁用转发 -->
        <el-table-column label="禁用转发消息" width="140">
          <template #default="{ row }">
            <el-switch
              :model-value="row.banForward === 1"
              @change="val => onBanForwardChange(val, row)"
            />
          </template>
        </el-table-column>

        <el-table-column label="禁用文本长度(-1无限制)">
          <template #default="{ row }">
            <el-input-number
              v-model="row.maxWords"
              :min="-1"
              @change="val => onMaxWordsChange(val, row)"
            />
            <div style="font-size:12px;color:#999">
              -1 表示不限制
            </div>
          </template>
        </el-table-column>

        

        <!-- <el-table-column label="创建时间" prop="createdAt" width="180" /> -->
      </el-table>

      <el-pagination
        layout="total, sizes, prev, pager, next"
        :total="total"
        :page-size="pageSize"
        :current-page="page"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import {
  getBotChatGroupList,
  updateBotChatGroup
} from '@/api/bot/botChatGroup'

import { getBotChoice } from '@/api/bot/bot'
import { ElMessage } from 'element-plus'
import { ref, onMounted } from 'vue'

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

/* 搜索 */
const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {}
  onSubmit()
}

/* 分页 */
const handleCurrentChange = p => {
  page.value = p
  getTableData()
}

const handleSizeChange = s => {
  pageSize.value = s
  getTableData()
}

/* 行内更新 */
const updateRow = async (row, patch) => {
  const payload = {
    ID: row.ID,
    botID: row.botID,
    chatGroupID: row.chatGroupID,
    chatGroupName: row.chatGroupName,
    syncMessage: row.syncMessage,
    banForward: row.banForward,
    maxWords: row.maxWords,
    ...patch
  }

  const res = await updateBotChatGroup(payload)
  if (res.code === 0) {
    Object.assign(row, patch)
    ElMessage.success('更新成功')
  } else {
    ElMessage.error(res.msg || '更新失败')
  }
}

/* 最大字数 */
const onMaxWordsChange = (val, row) => {
  updateRow(row, { maxWords: val })
}

const onSyncChange = (val, row) => {
  updateRow(row, { syncMessage: val ? 1 : 2 })
}

const onBanForwardChange = (val, row) => {
  updateRow(row, { banForward: val ? 1 : 2 })
}
</script>
