<template>
<div>

  <el-button type="primary" @click="openDialog()">新增分组</el-button>

  <!-- 表格 -->
  <el-table :data="tableData">

    <el-table-column label="分组名称" prop="title" />

    <!-- 群 -->
    <el-table-column label="群组">
      <template #default="{ row }">
        <el-tag
          v-for="g in getNames(row.chatGroups, chatGroupMapper)"
          :key="g.id"
          style="margin:2px"
        >
          {{ g.name }}
        </el-tag>

        <el-tag v-if="hasMore(row.chatGroups)">...</el-tag>
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


  <!-- 群选择（已改：机器人 + 群） -->
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


  <!-- 用户选择（已改：加确认按钮） -->
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

/* 加载 */
const getTableData = async () => {
  const res = await getBotChatGroupClassifyList({ page:1,pageSize:10 })
  if (res.code === 0) {
    tableData.value = res.data.list
    chatGroupMapper.value = res.data.chatGroupMapper || {}
    userMapper.value = res.data.userMapper || {}
  }
}
onMounted(getTableData)


// ===== 工具方法 =====
const getNames = (idsStr, mapper) => {
  if (!idsStr) return []
  const ids = idsStr.split(',')
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

  groupIDs.value = row?.chatGroups
    ? row.chatGroups.split(',').map(Number)
    : []

  userIDs.value = row?.permitUsers
    ? row.permitUsers.split(',').map(Number)
    : []
}


// ===== 删除（新增） =====
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


// ===== 群选择（已改） =====
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
    }
  })

  groupSelectVisible.value = false
}


// ===== 用户选择（已改） =====
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
  const res = await saveBotChatGroupClassify({
    ...form.value,
    chatGroups: groupIDs.value.join(','),
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