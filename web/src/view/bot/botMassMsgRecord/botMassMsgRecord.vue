<template>
  <div>
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
      <el-table :data="tableData" row-key="ID">
        <el-table-column label="机器人" width="120">
          <template #default="scope">
            {{ scope.row.botName }}
          </template>
        </el-table-column>

        <el-table-column label="群聊" width="120">
          <template #default="scope">
            {{ scope.row.chatGroupName }}
          </template>
        </el-table-column>

        <el-table-column label="发送消息" width="200">
          <template #default="scope">
            <el-button link type="primary" @click="openMsg(scope.row.msg)">
              查看
            </el-button>
          </template>
        </el-table-column>

        <el-table-column label="发送成员" prop="members" minWidth="300" />

        <el-table-column label="日期" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
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

    <!-- ================= 查看发送内容 ================= -->
    <el-drawer
      v-model="contentVisible"
      title="发送内容"
      size="50%"
      destroy-on-close
    >
      <RichView v-model="currentContent" />
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { formatDate } from '@/utils/format'
import RichView from '@/components/richtext/rich-view.vue'

import { getBotChoiceWithChatGroup } from '@/api/bot/bot'
import { massMsgHistory } from '@/api/bot/msgMass'

// ================= 查询条件 =================
const searchInfo = ref({
  botChatGroup: [], // 级联选择值 [botID, chatGroupID]
  createdAtRange: []
})

const elSearchFormRef = ref()

// ================= 机器人 & 群聊 =================
const botList = ref([])
const botCascaderOptions = ref([])

const loadBots = async () => {
  const res = await getBotChoiceWithChatGroup()
  if (res.code === 0) {
    botList.value = res.data
    // 转成级联选择
    botCascaderOptions.value = botList.value.map(bot => ({
      value: bot.botID,
      label: bot.name,
      children: (bot.botChatGroups || []).map(g => ({
        value: g.chatGroupID,
        label: g.chatGroupName
      }))
    }))
  }
}

// ================= 表格数据 =================
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const getTableData = async () => {
  let query = {
    page: page.value,
    pageSize: pageSize.value,
    createdAtRange: searchInfo.value.createdAtRange
  }

  if (searchInfo.value.botChatGroup.length === 2) {
    query.botID = searchInfo.value.botChatGroup[0]
    query.chatGroupID = searchInfo.value.botChatGroup[1]
  }

  const res = await massMsgHistory(query)
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

onMounted(() => {
  loadBots()
  getTableData()
})

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {
    botChatGroup: [],
    createdAtRange: []
  }
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// ================= 查看发送内容 =================
const contentVisible = ref(false)
const currentContent = ref('')

const openMsg = (msg) => {
  currentContent.value = msg
  contentVisible.value = true
}
</script>
