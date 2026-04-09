<template>
  <div>
    <!-- 表格 -->
    <div class="gva-table-box">
      <div style="margin-bottom: 10px;">
        <el-button type="primary" @click="openDialog()">新增分组</el-button>
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

    <el-table-column label="分组名称" prop="title" />

    <el-table-column label="群组列表">
      <template #default="{ row }">
        <el-tag
          v-for="g in getGroupNames(row)"
          :key="g.chatGroupID"
          style="margin:2px"
        >
          {{ g.chatGroupName }}
        </el-tag>
      </template>
    </el-table-column>

    <el-table-column label="操作" width="220">
      <template #default="{ row }">
        <el-button link type="primary" @click="onView(row)">查看</el-button>
        <el-button link type="primary" @click="openDialog(row)">编辑</el-button>
        <el-button link type="danger" @click="onDeleteRow(row)">删除</el-button>
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

<!-- 弹窗 -->
<el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
  <el-form label-width="80px">
    <el-form-item label="分组名称">
      <el-input v-model="form.title" placeholder="请输入分组名称" />
    </el-form-item>

    <el-form-item label="选择群组">
      <el-checkbox-group v-model="checkedGroups">
        <el-checkbox
          v-for="g in allGroups"
          :key="g.chatGroupID"
          :label="g.chatGroupID"
        >
          {{ g.chatGroupName }}
        </el-checkbox>
      </el-checkbox-group>
    </el-form-item>
  </el-form>

  <template #footer>
    <el-button @click="dialogVisible=false">取消</el-button>
    <el-button type="primary" @click="onSubmit">保存</el-button>
  </template>
</el-dialog>

<!-- 查看弹窗 -->
<el-dialog v-model="viewVisible" title="查看分组" width="400px">
  <p><b>分组名称：</b>{{ viewData.title }}</p>
  <div>
    <b>群组：</b>
    <el-tag
      v-for="g in getGroupNames(viewData)"
      :key="g.chatGroupID"
      style="margin:2px"
    >
      {{ g.chatGroupName }}
    </el-tag>
  </div>
</el-dialog>
  </div>
</template>

<script setup>
import {
  getBotChatGroupClassifyList,
  saveBotChatGroupClassify,
  deleteBotChatGroupClassify,
} from '@/api/bot/botChatGroup'

import { getBotChoice } from '@/api/bot/bot'

import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

/* 表格数据 */
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

/* 多选 */
const multipleSelection = ref([])

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

/* 获取列表 */
const getTableData = async () => {
  const res = await getBotChatGroupClassifyList({
    page: page.value,
    pageSize: pageSize.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

onMounted(getTableData)

/* 所有群 */
const allGroups = ref([])

const loadGroups = async () => {
  const res = await getBotChatGroupList({ page: 1, pageSize: 999 })
  if (res.code === 0) {
    allGroups.value = res.data.list
  }
}
loadGroups()

/* 显示群名称 */
const getGroupNames = (row) => {
  if (!row || !row.chatGroups) return []
  const ids = row.chatGroups.split(',').map(Number)
  return allGroups.value.filter(g => ids.includes(g.chatGroupID))
}

/* 弹窗 */
const dialogVisible = ref(false)
const dialogTitle = ref('新增分组')
const form = ref({})
const checkedGroups = ref([])

/* 打开弹窗 */
const openDialog = (row) => {
  dialogVisible.value = true

  if (row) {
    dialogTitle.value = '编辑分组'
    form.value = { ...row }
    checkedGroups.value = row.chatGroups
      ? row.chatGroups.split(',').map(Number)
      : []
  } else {
    dialogTitle.value = '新增分组'
    form.value = {}
    checkedGroups.value = []
  }
}

/* 保存 */
const onSubmit = async () => {
  const res = await saveBotChatGroupClassify({
    ...form.value,
    chatGroups: checkedGroups.value.join(',')
  })

  if (res.code === 0) {
    ElMessage.success('操作成功')
    dialogVisible.value = false
    getTableData()
  } else {
    ElMessage.error(res.msg || '失败')
  }
}

/* 删除（批量） */
const onDelete = async () => {
  if (!multipleSelection.value.length) {
    ElMessage.warning('请选择数据')
    return
  }

  await ElMessageBox.confirm('确定删除吗？', '提示')

  const ids = multipleSelection.value.map(i => i.ID)

  const res = await deleteBotChatGroupClassify({ ids })

  if (res.code === 0) {
    ElMessage.success('删除成功')
    getTableData()
  }
}

/* 删除（单个） */
const onDeleteRow = async (row) => {
  await ElMessageBox.confirm('确定删除该分组吗？', '提示')

  const res = await deleteBotChatGroupClassify({ ids: [row.ID] })

  if (res.code === 0) {
    ElMessage.success('删除成功')
    getTableData()
  }
}

/* 查看 */
const viewVisible = ref(false)
const viewData = ref({})

const onView = (row) => {
  viewData.value = row
  viewVisible.value = true
}

/* 分页 */
const handleCurrentChange = (p) => {
  page.value = p
  getTableData()
}
</script>
