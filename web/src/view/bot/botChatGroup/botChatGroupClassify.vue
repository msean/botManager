<template>
<div>

  <el-button type="primary" @click="openDialog()">新增分组</el-button>

  <!-- 表格 -->
  <el-table :data="tableData">

    <el-table-column label="分组名称" prop="title" />

    <!-- ✅ 群（只改这里） -->
    <el-table-column label="群组">
      <template #default="{ row }">

        <div
          v-for="item in getGroupedGroups(row.chatGroups)"
          :key="item.botID"
          style="margin-bottom:4px"
        >
          <!-- 机器人 -->
          <el-tag type="success" style="margin-right:4px">
            {{ botNameMap[item.botID] || item.botID }}
          </el-tag>

          <!-- 群 -->
          <el-tag
            v-for="gid in item.groups"
            :key="gid"
            style="margin:2px"
          >
            {{ chatGroupMapper[gid] || gid }}
          </el-tag>
        </div>

      </template>
    </el-table-column>

    <!-- 用户 -->
    <el-table-column label="用户">
      <template #default="{ row }">
        <el-tag
          v-for="u in getNames(row.permitUsers, userMapper)"
          :key="u.id"
          style="margin:2px"
        >
          {{ u.name }}
        </el-tag>

        <el-tag v-if="hasMore(row.permitUsers)">...</el-tag>
      </template>
    </el-table-column>

    <!-- 操作 -->
    <el-table-column label="操作">
      <template #default="{ row }">
        <el-button link @click="openDialog(row)">编辑</el-button>
        <el-button link type="danger" @click="onDelete(row)">删除</el-button>
      </template>
    </el-table-column>

  </el-table>


  <!-- 编辑弹窗 -->
  <el-dialog v-model="dialogVisible" title="编辑分组">

    <el-form>
      <el-form-item label="名称">
        <el-input v-model="form.title" />
      </el-form-item>

      <!-- 群 -->
      <el-form-item label="群组">
        <el-tag
          v-for="id in groupIDs"
          :key="id"
          closable
          @close="groupIDs = groupIDs.filter(i => i !== id)"
        >
          {{ chatGroupMapper[id] || id }}
        </el-tag>

        <el-button @click="openGroupSelect">➕</el-button>
      </el-form-item>

      <!-- 用户 -->
      <el-form-item label="用户">
        <el-tag
          v-for="id in userIDs"
          :key="id"
          closable
          @close="userIDs = userIDs.filter(i => i !== id)"
        >
          {{ userMapper[id] || id }}
        </el-tag>

        <el-button @click="openUserSelect">➕</el-button>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="dialogVisible=false">取消</el-button>
      <el-button type="primary" @click="onSubmit">保存</el-button>
    </template>

  </el-dialog>


  <!-- 群选择 -->
  <el-dialog v-model="groupSelectVisible" title="选择群" @closed="onGroupDialogClose">

    <el-select v-model="selectedBot" placeholder="选择机器人" style="width:100%">
      <el-option
        v-for="b in botList"
        :key="b.botID"
        :label="b.name"
        :value="b.botID"
      />
    </el-select>

    <el-checkbox-group v-model="tempGroupIDs" style="margin-top:10px">
      <el-checkbox
        v-for="g in currentGroups"
        :key="g.chatGroupID"
        :label="g.chatGroupID"
      >
        {{ g.chatGroupName }}
      </el-checkbox>
    </el-checkbox-group>

    <template #footer>
      <el-button @click="groupSelectVisible=false">取消</el-button>
      <el-button type="primary" @click="confirmGroupSelect">确认</el-button>
    </template>

  </el-dialog>


  <!-- 用户选择 -->
  <el-dialog v-model="userSelectVisible" title="选择用户">
    <el-checkbox-group v-model="userIDs">
      <el-checkbox
        v-for="u in allUsers"
        :key="u.ID"
        :label="u.ID"
      >
        {{ u.nickName }}
      </el-checkbox>
    </el-checkbox-group>

    <template #footer>
      <el-button @click="userSelectVisible=false">取消</el-button>
      <el-button type="primary" @click="confirmUserSelect">确认</el-button>
    </template>
  </el-dialog>

</div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getBotChatGroupClassifyList,
  saveBotChatGroupClassify,
  deleteBotChatGroupClassify
} from '@/api/bot/botChatGroup'

import { getBotChoiceWithChatGroup } from '@/api/bot/bot'
import { userAll } from '@/api/user'

/* 数据 */
const tableData = ref([])
const chatGroupMapper = ref({})
const userMapper = ref({})
const botNameMap = ref({}) // ✅ 新增

/* 群对应bot */
const groupBotMap = ref({})

/* 加载 */
const getTableData = async () => {
  const res = await getBotChatGroupClassifyList({ page:1,pageSize:10 })
  if (res.code === 0) {
    tableData.value = res.data.list
    chatGroupMapper.value = res.data.chatGroupMapper || {}
    userMapper.value = res.data.userMapper || {}
    botNameMap.value = res.data.botMapper || {} // ✅ 新增
  }
}
onMounted(getTableData)


// ===== ✅ 新增方法（核心）=====
const getGroupedGroups = (str) => {
  if (!str) return []

  const map = {}

  str.split(',').forEach(item => {
    const arr = item.split('_')
    if (arr.length === 2) {
      const botID = Number(arr[0])
      const groupID = Number(arr[1])

      if (!map[botID]) map[botID] = []
      map[botID].push(groupID)
    }
  })

  return Object.keys(map).map(botID => ({
    botID: Number(botID),
    groups: map[botID]
  }))
}


// ===== 原有代码不动 =====
const getNames = (idsStr, mapper) => {
  if (!idsStr) return []

  const ids = idsStr.split(',').map(item => {
    const arr = item.split('_')
    return arr.length === 2 ? arr[1] : item
  })

  return ids.slice(0,10).map(id => ({
    id,
    name: mapper[id] || id
  }))
}

const hasMore = (idsStr) => {
  if (!idsStr) return false
  return idsStr.split(',').length > 10
}


// ===== 编辑 =====
const dialogVisible = ref(false)
const form = ref({})
const groupIDs = ref([])
const userIDs = ref([])

const openDialog = (row) => {
  dialogVisible.value = true
  form.value = row || {}

  if (row?.chatGroups) {
    const arr = row.chatGroups.split(',')

    groupIDs.value = arr.map(item => {
      const parts = item.split('_')
      if (parts.length === 2) {
        const botID = Number(parts[0])
        const groupID = Number(parts[1])
        groupBotMap.value[groupID] = botID
        return groupID
      }
      return Number(item)
    })
  } else {
    groupIDs.value = []
  }

  userIDs.value = row?.permitUsers
    ? row.permitUsers.split(',').map(Number)
    : []
}


// ===== 删除 =====
const onDelete = async (row) => {
  await ElMessageBox.confirm('确认删除该分组？')

  const res = await deleteBotChatGroupClassify({
    ids: [row.ID]
  })

  if (res.code === 0) {
    ElMessage.success('删除成功')
    getTableData()
  }
}


// ===== 群选择 =====
const groupSelectVisible = ref(false)
const botList = ref([])
const selectedBot = ref(null)
const currentGroups = ref([])
const tempGroupIDs = ref([])

const openGroupSelect = async () => {
  groupSelectVisible.value = true

  selectedBot.value = null
  currentGroups.value = []

  tempGroupIDs.value = [...groupIDs.value]

  const res = await getBotChoiceWithChatGroup()
  botList.value = res.data || []
}

watch(selectedBot, (val) => {
  const bot = botList.value.find(b => b.botID === val)
  currentGroups.value = bot?.botChatGroups || []
})

const confirmGroupSelect = () => {
  const set = new Set([...groupIDs.value, ...tempGroupIDs.value])
  groupIDs.value = Array.from(set)

  currentGroups.value.forEach(g => {
    if (groupIDs.value.includes(g.chatGroupID)) {
      chatGroupMapper.value[g.chatGroupID] = g.chatGroupName
      groupBotMap.value[g.chatGroupID] = selectedBot.value
    }
  })

  groupSelectVisible.value = false
}


// ===== 用户选择 =====
const userSelectVisible = ref(false)
const allUsers = ref([])

const openUserSelect = async () => {
  userSelectVisible.value = true
  const res = await userAll()
  allUsers.value = res.data.list || []
}

const confirmUserSelect = () => {
  allUsers.value.forEach(u => {
    if (userIDs.value.includes(u.ID)) {
      userMapper.value[u.ID] = u.nickName
    }
  })
  userSelectVisible.value = false
}


// ===== 保存 =====
const onSubmit = async () => {
  const chatGroups = groupIDs.value.map(gid => {
    const botID = groupBotMap.value[gid] || 0
    return `${botID}_${gid}`
  })

  const res = await saveBotChatGroupClassify({
    ...form.value,
    chatGroups: chatGroups.join(','),
    permitUsers: userIDs.value.join(','),
    refresh: true
  })

  if (res.code === 0) {
    ElMessage.success('成功')
    dialogVisible.value = false
    getTableData()
  }
}

const onGroupDialogClose = () => {
  tempGroupIDs.value = []
  selectedBot.value = null
  currentGroups.value = []
}
</script>