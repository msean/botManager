<template>
  <div>
    <!-- 搜索 -->
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo">
        <el-form-item label="群 / 频道">
          <el-select
            v-model="searchInfo.groupId"
            placeholder="请选择群 / 频道"
            filterable
            clearable
            style="width:260px"
          >
            <el-option
              v-for="item in chatOptions"
              :key="item.groupId"
              :label="item.groupName"
              :value="item.groupId"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="关键词">
          <el-input
            v-model="searchInfo.keyword"
            clearable
            style="width:200px"
          />
        </el-form-item>

        <el-form-item label="时间">
          <el-date-picker
            v-model="searchInfo.timeRange"
            type="datetimerange"
            class="!w-380px"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
          <el-button type="success" @click="onExport">导出</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格 -->
    <div class="gva-table-box">
      <el-table :data="tableData" style="width:100%">
        <el-table-column label="群" prop="groupName" width="160" />

        <el-table-column label="用户" width="160">
          <template #default="scope">
            {{ scope.row.nickName || scope.row.username || scope.row.userId }}
          </template>
        </el-table-column>

        <el-table-column label="内容" prop="text" min-width="300" />
        <el-table-column label="类型" prop="messageType" width="100" />

        <el-table-column label="时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.timestamp) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next"
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { formatDate } from '@/utils/format'
import { getListenChoice, getListenList, exportListen } from '@/api/listen/listen'

defineOptions({ name: 'TelegramListen' })

/* ===== 搜索条件 ===== */
const searchInfo = ref({
  groupId: null,
  keyword: '',
  timeRange: []
})

/* ===== 群 / 频道 ===== */
const chatOptions = ref([])
const loadChatOptions = async () => {
  const res = await getListenChoice()
  if (res.code === 0) {
    chatOptions.value = res.data.map(i => ({
      groupId: i.group_id,
      groupName: i.group_name
    }))
  }
}

/* ===== 表格 ===== */
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const buildParams = () => {
  const params = {
    groupId: searchInfo.value.groupId,
    keyword: searchInfo.value.keyword || '',
    page: page.value,
    pageSize: pageSize.value
  }

  if (searchInfo.value.timeRange?.length === 2) {
    params.startTime = searchInfo.value.timeRange[0]
    params.endTime = searchInfo.value.timeRange[1]
  }

  return params
}

const getTableData = async () => {
  if (!searchInfo.value.groupId) return

  const res = await getListenList(buildParams())
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

/* ===== 事件 ===== */
const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = { groupId: null, keyword: '', timeRange: [] }
  tableData.value = []
  total.value = 0
}

const handleCurrentChange = val => {
  page.value = val
  getTableData()
}

const handleSizeChange = val => {
  pageSize.value = val
  page.value = 1
  getTableData()
}

/* ===== 导出 ===== */
const onExport = async () => {
  if (!searchInfo.value.groupId) {
    ElMessage.warning('请选择群 / 频道')
    return
  }

  const params = buildParams()
  delete params.page
  delete params.pageSize

  const res = await exportListen(params)
  if (res.code === 0 && res.data?.url) {
    window.open(res.data.url)
    ElMessage.success('开始下载')
  }
}

loadChatOptions()
</script>
