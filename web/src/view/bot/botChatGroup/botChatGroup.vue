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
      <div style="margin-bottom: 10px;">
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
        <el-table-column label="机器人" prop="botName" width="180" />
        <el-table-column label="群组名称" prop="chatGroupName" width="280" />

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

        <!-- 最大字数 -->
        <el-table-column label="禁用文本长度" width="220">
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

        <!-- 必须关注频道 -->
        <el-table-column label="必须关注频道" min-width="260" align="center">
          <template #default="{ row }">
            <div style="display:flex; flex-wrap:wrap; gap:6px; justify-content:center;">
              <el-tag
                v-for="ch in getChannelNamesWithObj(row)"
                :key="ch.channelID"
                closable
                type="info"
                @close="removeChannel(row, ch.channelID)"
              >
                {{ ch.channelName }}
              </el-tag>

              <el-button
                link
                type="primary"
                icon="Plus"
                @click="openBindChannel(row)"
              />
            </div>
          </template>
        </el-table-column>

         <!-- 邀请链接 -->
        <el-table-column label="渠道关注链接" min-width="300" align="center">
          <template #default="{ row }">
            <el-input
              v-model="row.invaidChannelFoldLink"
              placeholder="填写邀请链接"
              @change="val => onLinkChange(val, row)"
            />
            <div style="font-size:12px;color:#999">
              例如：https://t.me/xxxx
            </div>
          </template>
        </el-table-column>
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

    <!-- 绑定频道弹窗 -->
    <el-dialog
      v-model="bindDialogVisible"
      title="绑定必须关注的频道"
      width="420px"
    >
      <el-checkbox-group v-model="checkedChannelIds">
        <el-checkbox
          v-for="ch in currentBotChannels"
          :key="ch.channelID"
          :label="ch.channelID"
        >
          {{ ch.channelName }}
        </el-checkbox>
      </el-checkbox-group>

      <template #footer>
        <el-button @click="bindDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitBindChannels">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  getBotChatGroupList,
  updateBotChatGroup,
  deleteBotChatGroupByIds
} from '@/api/bot/botChatGroup'

import { getBotChoice } from '@/api/bot/bot'
import { getBotChoiceWithChatGroup } from '@/api/bot/bot'

import { ElMessage } from 'element-plus'
import { ref, onMounted, computed } from 'vue'

const multipleSelection = ref([])

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const onDelete = async () => {
  if (!multipleSelection.value.length) {
    ElMessage.warning('请选择要删除的数据')
    return
  }

  const ids = multipleSelection.value.map(item => item.ID)

  const res = await deleteBotChatGroupByIds({ ids })

  if (res.code === 0) {
    ElMessage.success('删除成功')
    getTableData()
  } else {
    ElMessage.error(res.msg || '删除失败')
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

/* 机器人 + 频道 */
const botChannelOptions = ref([])
const loadBotChannels = async () => {
  const res = await getBotChoiceWithChatGroup()
  if (res.code === 0) {
    botChannelOptions.value = res.data
  }
}
loadBotChannels()

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
    mustJoinChannels: row.mustJoinChannels,
    InvaidChannelFoldLink: row.InvaidChannelFoldLink,
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

/* 原逻辑 */
const onMaxWordsChange = (val, row) => {
  updateRow(row, { maxWords: val })
}

const onSyncChange = (val, row) => {
  updateRow(row, { syncMessage: val ? 1 : 2 })
}

const onBanForwardChange = (val, row) => {
  updateRow(row, { banForward: val ? 1 : 2 })
}

/* 新增：邀请链接更新 */
const onLinkChange = (val, row) => {
  updateRow(row, { InvaidChannelFoldLink: val })
}

/* ================== 绑定频道 ================== */
const bindDialogVisible = ref(false)
const currentRow = ref(null)
const checkedChannelIds = ref([])

const openBindChannel = (row) => {
  currentRow.value = row
  bindDialogVisible.value = true
  checkedChannelIds.value = row.mustJoinChannels
    ? row.mustJoinChannels.split(',').map(Number)
    : []
}

const currentBotChannels = computed(() => {
  if (!currentRow.value) return []
  const bot = botChannelOptions.value.find(
    b => b.botID === currentRow.value.botID
  )
  return bot?.botChannels || []
})

const getChannelNames = (row) => {
  if (!row.mustJoinChannels) return []
  const ids = row.mustJoinChannels.split(',').map(Number)
  const bot = botChannelOptions.value.find(b => b.botID === row.botID)
  if (!bot) return []
  return bot.botChannels
    .filter(ch => ids.includes(ch.channelID))
    .map(ch => ch.channelName)
}

/* 新增：获取对象数组，用于可删除标签 */
const getChannelNamesWithObj = (row) => {
  if (!row.mustJoinChannels) return []
  const ids = row.mustJoinChannels.split(',').map(Number)
  const bot = botChannelOptions.value.find(b => b.botID === row.botID)
  if (!bot) return []
  return bot.botChannels.filter(ch => ids.includes(ch.channelID))
}

/* 删除某个频道 */
const removeChannel = async (row, channelID) => {
  const ids = getChannelNamesWithObj(row)
    .map(ch => ch.channelID)
    .filter(id => id !== channelID)

  await updateRow(row, {
    mustJoinChannels: ids.join(',')
  })
}

/* 提交绑定 */
const submitBindChannels = async () => {
  await updateRow(currentRow.value, {
    mustJoinChannels: checkedChannelIds.value.join(',')
  })
  bindDialogVisible.value = false
}
</script>
