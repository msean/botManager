
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

       <!-- 用户名 -->
      <el-form-item label="用户名" prop="userName">
        <el-input
          v-model="searchInfo.userName"
          placeholder="请输入用户名"
          clearable
          class="!w-200px"
        />
      </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <!-- <div class="gva-btn-list">
            <el-button  type="primary" icon="plus" @click="openDialog()">新增</el-button>
            <el-button  icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
            
        </div> -->
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        
            <el-table-column align="left" label="机器人名称" prop="botName" width="150" />

            <el-table-column align="left" label="用户名称" prop="userName" width="180" />

            <el-table-column align="left" label="昵称" prop="fullName" width="200" />

            <el-table-column align="left" label="群名" prop="chatName" width="360" />

            <el-table-column align="left" label="封禁时长" prop="banDuration" width="90">
              <template #default="scope">
                {{ scope.row.banDuration }} 分钟
              </template>
            </el-table-column>
            <el-table-column align="left" label="备注" prop="reMark" width="300" />
            <el-table-column align="left" label="状态" prop="status" width="300">
              <template #default="scope">
                <el-tag v-if="scope.row.status === 1" type="danger">
                  封禁中
                </el-tag>
                <el-tag v-else type="success">
                  非封禁中
                </el-tag>
              </template>
            </el-table-column>

            <el-table-column
              align="left"
              label="禁用类型"
              prop="banType"
              width="120"
            >
              <template #default="{ row }">
                {{ banTypeMap[row.banType] || '未知' }}
              </template>
            </el-table-column>

            <el-table-column  align="center" label="发送消息" prop="msg" width="120">
              <template #default="scope">
                <el-button type="primary" link @click="viewMsg(scope.row)">
                  查看
                </el-button>
              </template>
            </el-table-column>
            <el-table-column sortable align="left" label="创建时间" prop="createdAt" width="180">
              <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
            </el-table-column>

            <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
              <template #default="scope">
                <el-button
                  v-if="scope.row.status === 1"
                  type="danger"
                  link
                  class="table-button"
                  @click="handleUnban(scope.row)"
                >
                  解禁
                </el-button>
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
        <el-dialog
          v-model="dialogVisible"
          title="消息内容"
          width="600px"
        >
          <div style="white-space: pre-wrap;">{{ currentMsg }}</div>
          <template #footer>
            <el-button @click="dialogVisible = false">关闭</el-button>
          </template>
        </el-dialog>
    </div>
  </div>
</template>

<script setup>
import {
  createBanRecord,
  deleteBanRecord,
  deleteBanRecordByIds,
  updateBanRecord,
  findBanRecord,
  getBanRecordList
} from '@/api/bot/banRecord'

// 全量引入格式化工具 请按需保留
import { formatDate} from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"
import { unBanUser } from '@/api/bot/bot'

defineOptions({
    name: 'BanRecord'
})

const dialogVisible = ref(false)
const currentMsg = ref('')

const viewMsg = (row) => {
  currentMsg.value = row.msg
  dialogVisible.value = true
}
// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
            botID: undefined,
            userID: undefined,
            userName: '',
            chatID: undefined,
            chatName: '',
            banDuration: undefined,
        })

const banTypeMap = {
  1: '消息',
  2: '成员',
  3: '转发' // 如果以后有第三种类型
}

const handleUnban = (row) => {
  ElMessageBox.confirm(
    `确定要解禁用户：${row.userName || row.userID} 吗？`,
    '解禁确认',
    {
      confirmButtonText: '确定解禁',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const res = await unBanUser({
        id: row.ID,
        botID: row.botID,
        userID: row.userID,
        chatID: row.chatID
      })

      if (res.code === 0) {
        ElMessage.success('解禁成功')
        getTableData() // 刷新表格
      } else {
        ElMessage.error(res.msg || '解禁失败')
      }
    } catch (err) {
      ElMessage.error('解禁请求异常')
    }
  }).catch(() => {})
}

// 验证规则
const rule = reactive({
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getBanRecordList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteBanRecordFunc(row)
        })
    }

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          IDs.push(item.ID)
        })
      const res = await deleteBanRecordByIds({ IDs })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === IDs.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
      })
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateBanRecordFunc = async(row) => {
    const res = await findBanRecord({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteBanRecordFunc = async (row) => {
    const res = await deleteBanRecord({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        botID: undefined,
        userID: undefined,
        userName: '',
        chatID: undefined,
        chatName: '',
        banDuration: undefined,
        }
}
// 弹窗确定
const enterDialog = async () => {
     btnLoading.value = true
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return btnLoading.value = false
              let res
              switch (type.value) {
                case 'create':
                  res = await createBanRecord(formData.value)
                  break
                case 'update':
                  res = await updateBanRecord(formData.value)
                  break
                default:
                  res = await createBanRecord(formData.value)
                  break
              }
              btnLoading.value = false
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
                })
                closeDialog()
                getTableData()
              }
      })
}

const detailForm = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findBanRecord({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}


</script>

<style>

</style>
