
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
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
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


        <el-table-column align="left" label="机器人" prop="botName" width="240" />
        <el-table-column align="left" label="命令" prop="cmd" width="150" />

        <el-table-column label="设置内容" prop="startContent" width="200">
          <template #default>[富文本内容]</template>
        </el-table-column>

        <el-table-column label="按钮" width="300">
          <template #default="scope">
            <div v-if="scope.row.cmdButtons">
              <!-- 尝试解析 JSON -->
              <div
                v-for="(row, rowIndex) in parseButtons(scope.row.cmdButtons)"
                :key="rowIndex"
                style="margin-bottom: 4px;"
              >
                <span
                  v-for="(btn, btnIndex) in row"
                  :key="btnIndex"
                  style="margin-right: 6px;"
                >
                  {{ btn.name }} ({{ btn.bindCmd }})
                </span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column sortable align="left" label="创建日期" prop="createdAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateBotCmdConfigFunc(scope.row)">编辑</el-button>
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

    <!-- ======================== 新增/编辑抽屉 ======================== -->
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type==='create'?'新增':'编辑' }}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule">
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

        <el-form-item label="命令" prop="cmd">
          <el-input v-model="formData.cmd" clearable placeholder="请输入命令" />
        </el-form-item>

        <el-form-item label="开始设置内容:" prop="content">
          <RichEdit v-model="formData.content" />
        </el-form-item>

        <el-form-item label="配置按钮">
          <div class="btn-group-wrapper">

            <!-- 每一行按钮 -->
            <div v-for="(row, rowIndex) in formData.cmdButtons" :key="rowIndex" class="btn-row">
              <div
                v-for="(btn, btnIndex) in row"
                :key="btnIndex"
                class="btn-item"
              >
                <span class="btn-text">{{ btn.name }}</span>
                <el-icon class="btn-edit" @click.stop="openEditDialog(rowIndex, btnIndex)"><Edit /></el-icon>
                <el-icon class="btn-delete" @click.stop="removeButton(rowIndex, btnIndex)"><Close /></el-icon>
              </div>

              <el-button class="add-btn" type="primary" link @click="openAddDialog(rowIndex)">+ 添加按钮</el-button>
              <el-button class="delete-row-btn" type="danger" link @click="removeRow(rowIndex)">删除该行</el-button>
            </div>

            <el-button type="primary" plain size="small" @click="addNewRow" style="margin-top: 10px;">+ 新增一行</el-button>
          </div>

          <!-- 编辑按钮弹窗 -->
          <el-dialog v-model="dialogVisible" title="编辑按钮" width="400px">
            <el-form :model="editForm" label-width="90px">
              <el-form-item label="名称"><el-input v-model="editForm.name" /></el-form-item>
              <el-form-item label="绑定命令"><el-input v-model="editForm.bindCmd" /></el-form-item>
            </el-form>

            <template #footer>
              <el-button @click="dialogVisible = false">取消</el-button>
              <el-button type="primary" @click="saveButton">保存</el-button>
            </template>
          </el-dialog>
        </el-form-item>

      </el-form>
    </el-drawer>

    <!-- ============================ 详情抽屉 ========================== -->
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" title="查看">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="机器人">{{ detailForm.botName }}</el-descriptions-item>
        <el-descriptions-item label="回复内容"><RichView v-model="detailForm.content" /></el-descriptions-item>
        <el-descriptions-item label="命令按钮配置">
          <div v-if="parseButtons(detailForm.cmdButtons).length">
            <div 
              v-for="(row, rowIndex) in parseButtons(detailForm.cmdButtons)" 
              :key="rowIndex" 
              style="margin-bottom: 6px;"
            >
              <el-tag
                v-for="(btn, btnIndex) in row"
                :key="btnIndex"
                style="margin-right: 6px;"
                type="info"
              >
                {{ btn.name }} ({{ btn.bindCmd }})
              </el-tag>
            </div>
          </div>
        <div v-else>无</div>
      </el-descriptions-item>

      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createBotCmdConfig,
  deleteBotCmdConfig,
  deleteBotCmdConfigByIds,
  updateBotCmdConfig,
  findBotCmdConfig,
  getBotCmdConfigList
} from '@/api/bot/botCmdConfig'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'
import RichView from '@/components/richtext/rich-view.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, onMounted, reactive } from 'vue'
import { useAppStore } from "@/pinia"
import { getBotChoice } from '@/api/bot/bot'

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  botID: undefined,
  startContent: '',
  cmdButtons: [],          // ⭐ array，而不是字符串
})

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
  const table = await getBotCmdConfigList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
            deleteBotCmdConfigFunc(row)
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
      const res = await deleteBotCmdConfigByIds({ IDs })
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

const parseButtons = (buttons) => {
  try {
    return JSON.parse(buttons || '[]')
  } catch (e) {
    return []
  }
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateBotCmdConfigFunc = async (row) => {
  const res = await findBotCmdConfig({ ID: row.ID })
  if (res.code === 0) {
    let cmdButtonsData = res.data.cmdButtons || []
    if (typeof cmdButtonsData === 'string') {
      try {
        cmdButtonsData = JSON.parse(cmdButtonsData)
      } catch (err) {
        cmdButtonsData = []
      }
    }

    formData.value = {
      ...res.data,
      cmdButtons: cmdButtonsData
    }

    type.value = 'update'
    dialogFormVisible.value = true
  }
}



// 删除行
const deleteBotCmdConfigFunc = async (row) => {
    const res = await deleteBotCmdConfig({ ID: row.ID })
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
    startContent: '',
    cmdButtons: [],        // ⭐ 这里也要是数组
  }
}

// ======= 获取机器人列表 =======
const botList = ref([])
const getBotList = async () => {
  const res = await getBotChoice()
  if (res.code === 0) {
    botList.value = res.data || []
  }
}

// 弹窗确定
const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return (btnLoading.value = false)

    const payload = {
      ...formData.value,
      cmdButtons: formData.value.cmdButtons, 
    }

    let res
    if (type.value === 'create') res = await createBotCmdConfig(payload)
    else res = await updateBotCmdConfig(payload)

    btnLoading.value = false

    if (res.code === 0) {
      ElMessage.success("创建/更新成功")
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
  const res = await findBotCmdConfig({ ID: row.ID })
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


defineOptions({ name: 'BotCmdConfig' })

/* ------------ 按钮编辑控制 ------------ */

const dialogVisible = ref(false)
const editForm = ref({ name: '', bindCmd: '' })
let editRow = -1
let editIndex = -1

const openAddDialog = (rowIndex) => {
  editRow = rowIndex
  editIndex = -1
  editForm.value = { name: '', bindCmd: '' }
  dialogVisible.value = true
}

const openEditDialog = (rowIndex, btnIndex) => {
  editRow = rowIndex
  editIndex = btnIndex
  editForm.value = { ...formData.value.cmdButtons[rowIndex][btnIndex] }
  dialogVisible.value = true
}

const saveButton = () => {
  if (editIndex === -1) {
    // 新增按钮 → 推入当前行 editRow 对应的数组
    if (!formData.value.cmdButtons[editRow]) {
      formData.value.cmdButtons[editRow] = [] // 防止当前行不存在
    }
    formData.value.cmdButtons[editRow].push({ ...editForm.value })
  } else {
    // 编辑已有按钮
    formData.value.cmdButtons[editRow][editIndex] = { ...editForm.value }
  }
  dialogVisible.value = false
}


const removeButton = (rowIndex, btnIndex) => {
  formData.value.cmdButtons[rowIndex].splice(btnIndex, 1)
}

const removeRow = (rowIndex) => {
  formData.value.cmdButtons.splice(rowIndex, 1)
}

const addNewRow = () => {
  formData.value.cmdButtons.push([]) // 新增一行空按钮
}

onMounted(() => {
  getBotList()
})



</script>

<style>

.btn-row {
  display: flex;         /* 水平排列一行 */
  flex-wrap: nowrap;     /* 不换行 */
  align-items: center;   /* 垂直居中 */
  gap: 10px;             /* 按钮间距 */
}

.btn-item {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  background-color: #f5f5f5;
  border-radius: 4px;
  white-space: nowrap;   /* 防止文字换行 */
}

.add-btn,
.delete-row-btn {
  flex-shrink: 0;        /* 按钮保持原大小 */
}


</style>
