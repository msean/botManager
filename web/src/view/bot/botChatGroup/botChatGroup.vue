<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="创建日期" prop="createdAtRange">
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>

          <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>

         <!-- 机器人选择下拉框 -->
        <el-form-item label="机器人" prop="botID">
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

        <template v-if="showAllQuery">
          <!-- 可扩展更多查询条件 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <!-- <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true" v-if="!showAllQuery">展开</el-button> -->
          <!-- <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>收起</el-button> -->
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">新增</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column
          align="left"
          label="机器人名称"
          prop="botName"
          width="240"
        />

        <el-table-column
          align="left"
          label="群组名称"
          prop="chatGroupName"
          width="400"
        />

        <el-table-column
          align="left"
          label="开启消息同步"
          width="200"
        >
          <template #default="{ row }">
            <el-switch
              :model-value="row.syncMessage === 1"
              active-text="开启"
              inactive-text="关闭"
              @change="val => onSyncChange(val, row)"
            />
          </template>
        </el-table-column>
        <el-table-column
          sortable
          align="left"
          label="日期"
          prop="createdAt"
          width="180"
        >
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增' : '编辑' }}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="机器人:" prop="botID">
          <el-select v-model="formData.botID" placeholder="请选择机器人" clearable filterable style="width: 100%">
            <el-option
              v-for="bot in botList"
              :key="bot.botID"
              :label="bot.name"
              :value="bot.botID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="群组ID:" prop="chatGroupID">
          <el-input v-model.number="formData.chatGroupID" clearable placeholder="请输入群组ID" />
        </el-form-item>

        <el-form-item label="群组名称:" prop="chatGroupName">
          <el-input v-model="formData.chatGroupName" clearable placeholder="请输入群组名称" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createBotChatGroup,
  deleteBotChatGroup,
  deleteBotChatGroupByIds,
  updateBotChatGroup,
  findBotChatGroup,
  getBotChatGroupList
} from '@/api/bot/botChatGroup'

import { getBotChoice } from '@/api/bot/bot'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'
import { useAppStore } from "@/pinia"

defineOptions({
  name: 'BotChatGroup'
})

const botOptions = ref([])

const setOptions = async () => {
  const res = await getBotChoice()
  if (res.code === 0) {
    botOptions.value = res.data.map(item => ({
      label: item.name,
      value: item.botID
    }))
  }
}
setOptions()

const appStore = useAppStore()
const btnLoading = ref(false)
const showAllQuery = ref(false)
const formData = ref({
  botID: undefined,
  chatGroupID: undefined,
  chatGroupName: ''
})

const rule = reactive({
  botID: [{ required: true, message: '请选择机器人', trigger: 'change' }],
  chatGroupID: [{ required: true, message: '请输入群组ID', trigger: 'blur' }],
  chatGroupName: [{ required: true, message: '请输入群组名称', trigger: 'blur' }]
})

const elFormRef = ref()
const elSearchFormRef = ref()

// 表格控制部分
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const multipleSelection = ref([])

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}



const onSubmit = () => {
  elSearchFormRef.value?.validate(async valid => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

const handleSizeChange = val => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = val => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const table = await getBotChatGroupList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}


// 🔹 获取机器人列表
const botList = ref([])
const loadBotList = async () => {
  const res = await getBotChoice()
  if (res.code === 0) {
    botList.value = res.data || []
  }
}

// 多选删除
const handleSelectionChange = val => {
  multipleSelection.value = val
}

const deleteRow = row => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => deleteBotChatGroupFunc(row))
}

const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const IDs = multipleSelection.value.map(i => i.ID)
    if (IDs.length === 0) {
      return ElMessage.warning('请选择要删除的数据')
    }
    const res = await deleteBotChatGroupByIds({ IDs })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === IDs.length && page.value > 1) page.value--
      getTableData()
    }
  })
}

const type = ref('')
const dialogFormVisible = ref(false)
const detailShow = ref(false)
const detailForm = ref({})

const updateBotChatGroupFunc = async row => {
  const res = await findBotChatGroup({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
    loadBotList()
  }
}

const deleteBotChatGroupFunc = async row => {
  const res = await deleteBotChatGroup({ ID: row.ID })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    if (tableData.value.length === 1 && page.value > 1) page.value--
    getTableData()
  }
}

const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
  loadBotList()
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = { botID: undefined, chatGroupID: undefined, chatGroupName: '' }
}

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async valid => {
    if (!valid) return (btnLoading.value = false)
    let res
    if (type.value === 'create') res = await createBotChatGroup(formData.value)
    else res = await updateBotChatGroup(formData.value)

    btnLoading.value = false
    if (res.code === 0) {
      ElMessage.success('创建/更改成功')
      closeDialog()
      getTableData()
    }
  })
}

const onSyncChange = async (val, row) => {
  // val: true / false
  // row: 当前整行数据

  const newSyncValue = val ? 1 : 2

  // 组装完整更新参数（整行一起更新）
  const payload = {
    ID: row.ID,
    botID: row.botID,
    chatGroupID: row.chatGroupID,
    chatGroupName: row.chatGroupName,
    syncMessage: newSyncValue
  }

  try {
    const res = await updateBotChatGroup(payload)
    if (res.code === 0) {
      ElMessage.success(newSyncValue === 1 ? '消息同步已开启' : '消息同步已关闭')
      // 直接更新当前行，避免重新拉列表
      row.syncMessage = newSyncValue
    } else {
      throw new Error(res.msg || '更新失败')
    }
  } catch (e) {
    ElMessage.error('更新失败，已回滚')
    // 回滚 UI
    row.syncMessage = row.syncMessage === 1 ? 2 : 1
  }
}



onMounted(() => {
  getTableData()
})
</script>

<style scoped></style>
