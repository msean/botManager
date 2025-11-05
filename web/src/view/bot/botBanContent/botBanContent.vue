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

        <el-form-item label="机器人" prop="botID">
          <el-select
            v-model="searchInfo.botID"
            placeholder="请选择机器人"
            filterable
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="item in botList"
              :key="item.botID"
              :label="item.name"
              :value="item.botID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="禁用内容" prop="banContent">
          <el-input v-model="searchInfo.banContent" placeholder="搜索条件" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
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
        <el-table-column align="left" label="机器人" prop="botName" width="180" />
        <el-table-column align="left" label="禁用内容" prop="banContent" width="300" />
        <el-table-column sortable align="left" label="创建日期" prop="createdAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateBotBanContentFunc(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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

    <!-- 抽屉 -->
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
          <el-select
            v-model="formData.botID"
            placeholder="请选择机器人"
            filterable
            clearable
            style="width: 100%"
            :disabled="isEdit"
          >
            <el-option
              v-for="item in botList"
              :key="item.botID"
              :label="item.name"
              :value="item.botID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="禁用内容:" prop="banContent">
          <el-input v-model="formData.banContent" clearable placeholder="请输入禁用内容" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createBotBanContent,
  deleteBotBanContent,
  deleteBotBanContentByIds,
  updateBotBanContent,
  findBotBanContent,
  getBotBanContentList
} from '@/api/bot/botBanContent'
import { getBotChoice } from '@/api/bot/bot'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted, computed } from 'vue'
import { useAppStore } from "@/pinia"
import { formatDate } from '@/utils/format'

defineOptions({ name: 'BotBanContent' })

const appStore = useAppStore()
const btnLoading = ref(false)
const showAllQuery = ref(false)
const dialogFormVisible = ref(false)
const elFormRef = ref()
const elSearchFormRef = ref()
const type = ref('')
const multipleSelection = ref([])

const formData = ref({
  banContent: '',
  botID: undefined
})
const botList = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const rule = reactive({
  banContent: [{ required: true, message: '请输入禁用内容', trigger: 'blur' }],
  botID: [{ required: true, message: '请选择机器人', trigger: 'change' }]
})

// ✅ 计算属性：是否编辑状态
const isEdit = computed(() => type.value === 'update')

// ======= 获取表格数据 =======
const getTableData = async () => {
  const res = await getBotBanContentList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

// ======= 搜索、重置 =======
const onSubmit = () => {
  page.value = 1
  getTableData()
}
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// ======= 分页 =======
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// ======= 获取机器人列表 =======
const getBotList = async () => {
  const res = await getBotChoice()
  if (res.code === 0) {
    botList.value = res.data || []
  }
}

// ======= 打开弹窗（新增） =======
const openDialog = () => {
  type.value = 'create'
  formData.value = {
    banContent: '',
    botID: undefined
  }
  dialogFormVisible.value = true
}

// ======= 编辑 =======
const updateBotBanContentFunc = async (row) => {
  const res = await findBotBanContent({ ID: row.ID })
  if (res.code === 0) {
    type.value = 'update'
    formData.value = res.data
    dialogFormVisible.value = true
  }
}

// ======= 保存 =======
const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return (btnLoading.value = false)
    let res
    if (type.value === 'update') {
      res = await updateBotBanContent(formData.value)
    } else {
      res = await createBotBanContent(formData.value)
    }
    btnLoading.value = false
    if (res.code === 0) {
      ElMessage.success('保存成功')
      dialogFormVisible.value = false
      getTableData()
    }
  })
}

// ======= 删除 =======
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗？', '提示', { type: 'warning' }).then(() => {
    deleteBotBanContentFunc(row)
  })
}
const deleteBotBanContentFunc = async (row) => {
  const res = await deleteBotBanContent({ ID: row.ID })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    if (tableData.value.length === 1 && page.value > 1) page.value--
    getTableData()
  }
}

// ======= 批量删除 =======
const onDelete = async () => {
  if (!multipleSelection.value.length) return ElMessage.warning('请选择要删除的数据')
  ElMessageBox.confirm('确定要删除选中的数据吗？', '提示', { type: 'warning' }).then(async () => {
    const IDs = multipleSelection.value.map(i => i.ID)
    const res = await deleteBotBanContentByIds({ IDs })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === IDs.length && page.value > 1) page.value--
      getTableData()
    }
  })
}

// ======= 多选 =======
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

onMounted(() => {
  getTableData()
  getBotList()
})
</script>
