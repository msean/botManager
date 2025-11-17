<template>
  <div>
    <!-- ================= 搜索表单 ================= -->
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
            clearable 
            style="width:200px"
          >
            <el-option 
              v-for="bot in botOptions" 
              :key="bot.botID" 
              :label="bot.name" 
              :value="bot.botID" 
            />
          </el-select>
        </el-form-item>

        <template v-if="showAllQuery">
          <!-- 展开更多查询项 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- ================= 表格 ================= -->
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
        <el-table-column sortable align="left" label="日期" prop="createdAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="标题" prop="title" width="120" />
        <el-table-column align="left" label="群聊名称" prop="chatGroupName" width="120" />
        <el-table-column align="left" label="机器人名称" prop="botName" width="120" />
        <el-table-column align="left" label="发送类型" prop="taskSendType" width="120">
          <template #default="scope">{{ renderSendType(scope.row.taskSendType) }}</template>
        </el-table-column>
        <el-table-column label="发送内容" prop="content" width="200">
          <template #default="scope">
            <template v-if="scope.row.taskSendType === 2">{{ scope.row.content }}</template>
            <template v-else-if="scope.row.taskSendType === 3 || scope.row.taskSendType === 4">
              <div v-for="(url, idx) in parseContent(scope.row.content)" :key="idx">
                <img v-if="scope.row.taskSendType===3" :src="url" style="max-width:100px; margin-right:5px;" />
                <video v-else controls style="max-width:150px; margin-right:5px;">
                  <source :src="url" type="video/mp4" />
                </video>
              </div>
            </template>
            <template v-else>[富文本内容]</template>
          </template>
        </el-table-column>
        <el-table-column label="扩展按钮" prop="extrendButton" width="300">
          <template #default="scope">
            <div v-if="scope.row.extrendButton && scope.row.extrendButton.length">
              <div 
                v-for="(row, rowIndex) in scope.row.extrendButton" 
                :key="rowIndex" 
                style="margin-bottom: 4px;"
              >
                <span 
                  v-for="(btn, btnIndex) in row" 
                  :key="btnIndex"
                  style="margin-right: 6px;"
                >
                  {{ btn.name }}
                </span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="发送间隔" prop="sendInterval" width="120" />
        <el-table-column align="left" label="下一次发送时间" prop="nextSendTime" width="180">
          <template #default="scope">{{ formatDate(scope.row.nextSendTime) }}</template>
        </el-table-column>
        <el-table-column align="left" label="上一次发送时间" prop="preSendTime" width="180">
          <template #default="scope">{{ formatDate(scope.row.preSendTime) }}</template>
        </el-table-column>
        <el-table-column align="left" label="状态" prop="status" width="120">
          <template #default="scope">
            {{ scope.row.status === 1 ? '运行' : scope.row.status === 2 ? '停止' : '' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateBotTaskFunc(scope.row)">编辑</el-button>
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

    <!-- ================= 新增 / 编辑抽屉 ================= -->
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

      <el-form :model="formData" :rules="rules" label-width="120px" ref="elFormRef" label-position="left" class="drawer-form">
          <el-col :span="16">
             <el-form-item label="标题" prop="title">
                <el-input v-model="formData.title" clearable placeholder="请输入标题" />
              </el-form-item>
          </el-col>
          <el-col :span="16">
            <el-form-item label="机器人" prop="botID">
              <el-select v-model="formData.botID" placeholder="请选择机器人" @change="handleBotChange" style="width:100%">
                <el-option v-for="bot in botOptions" :key="bot.botID" :label="bot.name" :value="bot.botID" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="16">
            <el-form-item label="群聊" prop="chatGroupID">
              <el-select v-model="formData.chatGroupID" placeholder="请选择群聊" style="width:100%">
                <el-option v-for="group in groupOptions" :key="group.chatGroupID" :label="group.chatGroupName" :value="group.chatGroupID" />
              </el-select>
            </el-form-item>
          </el-col>

        <!-- 发送类型 -->
        <el-form-item label="发送类型" prop="taskSendType">
          <el-radio-group v-model="formData.taskSendType" @change="handleSendTypeChange">
            <el-radio-button :label="1">预设页面</el-radio-button>
            <el-radio-button :label="2">文本</el-radio-button>
            <el-radio-button :label="3">图片</el-radio-button>
            <el-radio-button :label="4">视频</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <!-- 根据类型显示不同输入 -->
        <el-form-item label="发送内容" v-if="formData.taskSendType === 1">
          <RichEdit v-model="formData.content" />
        </el-form-item>
        <el-form-item label="发送内容" v-if="formData.taskSendType === 2">
          <el-input type="textarea" :rows="5" v-model="formData.content" placeholder="请输入文本内容" />
        </el-form-item>
        <el-form-item label="发送内容" v-if="formData.taskSendType === 3 || formData.taskSendType === 4">
          <el-upload
            :action="uploadUrl"
            list-type="picture-card"
            multiple
            :file-list="formData.uploadFileList"
            :on-success="handleUploadSuccess"
          >
            <i class="el-icon-plus"></i>
          </el-upload>
        </el-form-item>

       <!-- 扩展按钮 -->
      <el-form-item label="扩展按钮">
        <div class="btn-group-wrapper">
  <!-- 每行按钮 -->
  <div
    v-for="(row, rowIndex) in formData.extrendButton"
    :key="rowIndex"
    class="btn-row"
  >
    <!-- 行内按钮 -->
    <div
      v-for="(btn, btnIndex) in row"
      :key="btnIndex"
      class="btn-item"
    >
      <span class="btn-text">{{ btn.name }}</span>
      <el-icon class="btn-edit" @click.stop="openEditDialog(rowIndex, btnIndex)">
        <Edit />
      </el-icon>
      <el-icon class="btn-delete" @click.stop="removeButton(rowIndex, btnIndex)">
        <Close />
      </el-icon>
    </div>

    <!-- 当前行新增按钮 -->
    <el-button class="add-btn" type="primary" link @click="openAddDialog(rowIndex)">
      + 添加按钮
    </el-button>

    <!-- 删除整行 -->
    <el-button class="delete-row-btn" type="danger" link @click="removeRow(rowIndex)">
      删除该行
    </el-button>
  </div>

  <!-- 新增空行按钮 -->
  <el-button type="primary" plain size="small" @click="addNewRow" style="margin-top: 10px;">
    + 新增一行
  </el-button>
</div>



  <!-- 编辑/新增按钮 弹窗 -->
  <el-dialog v-model="dialogVisible" title="编辑按钮" width="400px">
    <el-form :model="editForm" label-width="90px">
      <el-form-item label="名称">
        <el-input v-model="editForm.name" />
      </el-form-item>
      <el-form-item label="跳转链接">
        <el-input v-model="editForm.url" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="saveButton">保存</el-button>
    </template>
  </el-dialog>
</el-form-item>


        <!-- 发送间隔 -->
        <el-col :span="16">
          <el-form-item label="发送间隔" prop="sendInterval">
            <el-input v-model.number="formData.sendInterval" placeholder="请输入发送间隔(分钟)" />
          </el-form-item>
        </el-col>

        <!-- 下一次发送时间 -->
        <el-col :span="16">
          <el-form-item label="下一次发送时间" prop="nextSendTimeStr">
            <el-date-picker
              v-model="formData.nextSendTimeStr"
              type="datetime"
              style="width: 100%"
              placeholder="选择日期时间"
              format="YYYY-MM-DD HH:mm"
              value-format="YYYY-MM-DD HH:mm:ss"
            />
          </el-form-item>
        </el-col>

        <!-- 下一次发送时间 -->
        <el-col :span="16">
          <el-form-item label="结束时间" prop="stopTimeText">
            <el-date-picker
              v-model="formData.stopTimeText"
              type="datetime"
              style="width: 100%"
              placeholder="选择日期时间"
              format="YYYY-MM-DD HH:mm"
              value-format="YYYY-MM-DD HH:mm:ss"
            />
          </el-form-item>
        </el-col>
        <!-- 上一次发送时间：仅编辑显示且禁用 -->
        <el-form-item label="上一次发送时间" v-if="type==='update'">
          <el-date-picker
            v-model="formData.preSendTime"
            type="datetime"
            style="width: 100%"
            placeholder="选择日期时间"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DD HH:mm:ss"
            disabled
          />
        </el-form-item>

        <!-- 状态开关 -->
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" active-text="开启" inactive-text="关闭" />
        </el-form-item>
      </el-form>


      <!-- 添加/编辑按钮弹窗 -->
      <el-dialog v-model="btnDialogVisible" :title="isEdit ? '编辑按钮' : '新增按钮'" width="400px">
        <el-form :model="btnForm" label-width="80px">
          <el-form-item label="名称">
            <el-input v-model="btnForm.name" />
          </el-form-item>
          <el-form-item label="链接">
            <el-input v-model="btnForm.url" />
          </el-form-item>
        </el-form>

        <template #footer>
          <el-button @click="btnDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmButton">确定</el-button>
        </template>
      </el-dialog>
    </el-drawer>

    <!-- ================= 详情抽屉 ================= -->
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
      <el-descriptions :column="1" border label-width="150px">
        <el-descriptions-item label="机器人名称">{{ detailForm.botName }}</el-descriptions-item>
        <el-descriptions-item label="群聊名称">{{ detailForm.chatGroupName }}</el-descriptions-item>
        <el-descriptions-item label="发送类型">{{ renderSendType(detailForm.taskSendType) }}</el-descriptions-item>
        <el-descriptions-item label="发送内容">
          <!-- 文本类型 -->
          <template v-if="detailForm.taskSendType === 2">
            {{ detailForm.content }}
          </template>

          <!-- 图片或视频类型 -->
          <template v-else-if="detailForm.taskSendType === 3 || detailForm.taskSendType === 4">
            <div v-for="(url, idx) in parseContent(detailForm.content)" :key="idx" style="margin-bottom:5px;">
              <img v-if="detailForm.taskSendType === 3" :src="url" style="max-width:100px; margin-right:5px;" />
              <video v-else controls style="max-width:150px; margin-right:5px;">
                <source :src="url" type="video/mp4" />
              </video>
            </div>
          </template>
          <!-- 富文本类型 -->
          <template v-else-if="detailForm.taskSendType === 1">
            <div v-html="safeHtml" class="rich-view-content"></div>
          </template>
        </el-descriptions-item>
        <el-descriptions-item label="扩展按钮">
          <div v-if="detailForm.extrendButton && detailForm.extrendButton.length">
            <div 
              v-for="(row, rowIndex) in detailForm.extrendButton" 
              :key="rowIndex" 
              style="margin-bottom: 4px;"
            >
              <span 
                v-for="(btn, btnIndex) in row" 
                :key="btnIndex"
                style="margin-right: 6px; padding: 2px 6px; background-color: #f5f7fa; border-radius: 4px;"
              >
                {{ btn.name }}
              </span>
            </div>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="发送间隔">{{ detailForm.sendInterval }}</el-descriptions-item>
        <el-descriptions-item label="下一次发送时间">{{ detailForm.nextSendTime }}</el-descriptions-item>
        <el-descriptions-item label="上一次发送时间">{{ detailForm.preSendTime }}</el-descriptions-item>
        <el-descriptions-item label="结束时间">{{ detailForm.stopTime }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          {{ detailForm.status === 1 ? '运行' : detailForm.status === 2 ? '停止' : '' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createBotTask, updateBotTask, deleteBotTaskByIds, findBotTask, getBotTaskList, deleteBotTask } from '@/api/bot/botTask'
import { getBotChoiceWithChatGroup } from '@/api/bot/bot'
import RichEdit from '@/components/richtext/rich-edit.vue'
import RichView from '@/components/richtext/rich-view.vue'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/pinia'
import { computed } from 'vue'
import DOMPurify from 'dompurify'
import ButtonGroupEditor from "@/components/ButtonGroupEditor.vue";


const appStore = useAppStore()
const btnLoading = ref(false)
const showAllQuery = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const searchInfo = ref({})
const multipleSelection = ref([])
const safeHtml = computed(() => {
  return DOMPurify.sanitize(detailForm.value?.content || '')
})

// ================== 表单数据 ==================
const formData = reactive({
  botID: undefined,
  chatGroupID: undefined,
  taskSendType: undefined,
  content: '',
  extrendButton: [],
  sendInterval: undefined,
  nextSendTime: new Date(),
  preSendTime: new Date(),
  status: true, // ⚠ 表单内部用 true/false
  uploadFileList: []
})

const rules = {
  nextSendTimeStr: [
    { required: true, message: '请选择下一次发送时间', trigger: 'change' }
  ],
  stopTimeText: [
    { required: true, message: '请选择结束时间', trigger: 'change' }
  ],
  sendInterval: [
    { required: true, message: '发送间隔不能为空', trigger: 'blur' }
  ]
}

const elFormRef = ref()
const elSearchFormRef = ref()

// ================== Bot + Group ==================
const botOptions = ref([])
const groupOptions = ref([])

const loadBotOptions = async () => {
  const res = await getBotChoiceWithChatGroup()
  if (res.code === 0) botOptions.value = res.data || []
}
loadBotOptions()

const handleBotChange = (botID) => {
  const bot = botOptions.value.find(b => b.botID === botID)
  groupOptions.value = bot ? bot.botChatGroups || [] : []
  formData.chatGroupID = undefined
}

// ================== 扩展按钮 ==================
const addExtraButton = () => formData.extrendButton.push({ name: '', url: '' })
const removeExtraButton = (index) => formData.extrendButton.splice(index, 1)
const uploadUrl = `${import.meta.env.VITE_BASE_API}/public/uploadMedia`
// ================== 上传 ==================
const handleUploadSuccess = (res, file) => {
  if (!res.url) {
    ElMessage.error("上传失败：后端未返回 url")
    return
  }

  // 如果 content 为空字符串，则初始化为数组
  let contentArray = formData.content ? JSON.parse(formData.content) : []
  contentArray.push(res.url)

  // 再把数组转成字符串赋值给 content
  formData.content = JSON.stringify(contentArray)
}

// ================== 表单重置 ==================
const resetFormData = () => {
  Object.assign(formData, {
    botID: undefined,
    chatGroupID: undefined,
    taskSendType: undefined,
    content: '',
    extrendButton: [],
    sendInterval: undefined,
    nextSendTime: new Date(),
    preSendTime: new Date(),
    status: true,  
    uploadFileList: [],
    nextSendTimeStr: '',
    stopTimeText: '',
  })
  groupOptions.value = []
}

// ================== 表格查询 ==================
const getTableData = async () => {
  const res = await getBotTaskList({
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
getTableData()

const onReset = () => { searchInfo.value = {}; getTableData() }
const onSubmit = () => { elSearchFormRef.value?.validate(v => v && (page.value = 1, getTableData())) }
const handleSizeChange = (val) => { page.value = 1; pageSize.value = val; getTableData() }
const handleCurrentChange = (val) => { page.value = val; getTableData() }
const handleSelectionChange = (val) => { multipleSelection.value = val }

const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { type: 'warning' })
    .then(() => deleteBotTaskFunc(row))
}

const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { type: 'warning' })
    .then(async () => {
      const IDs = multipleSelection.value.map(i => i.ID)
      if (!IDs.length) return ElMessage.warning('请选择要删除的数据')

      const res = await deleteBotTaskByIds({ IDs })
      if (res.code === 0) {
        ElMessage.success('删除成功')
        if (tableData.value.length === IDs.length && page.value > 1) page.value--
        getTableData()
      }
    })
}

// ================= 发送类型切换 =================
const handleSendTypeChange = (val) => {
  formData.taskSendType = Number(val)
  formData.content = ''
  formData.uploadFileList = []
}

// ================= 新增/编辑 Dialog =================
let type = ref('')
const dialogFormVisible = ref(false)

const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
  resetFormData()
}

const closeDialog = () => {
  dialogFormVisible.value = false
  resetFormData()
}

// ================== 保存按钮 ==================
const enterDialog = async () => {
  btnLoading.value = true

  elFormRef.value?.validate(async valid => {
    if (!valid) return btnLoading.value = false

    const payload = {
      ...formData,
      status: formData.status ? 1 : 2
    }

    let res
    if (type.value === 'create') res = await createBotTask(payload)
    else res = await updateBotTask(payload)

    btnLoading.value = false

    if (res.code === 0) {
      ElMessage.success('操作成功')
      closeDialog()
      getTableData()
    }
  })
}

// ================== 编辑 ==================
const updateBotTaskFunc = async (row) => {
  const res = await findBotTask({ ID: row.ID })
  if (res.code !== 0) return

  type.value = 'update'
  resetFormData()

  const data = res.data
  Object.assign(formData, {
    ...data,
    taskSendType: Number(data.taskSendType),
    status: data.status === 1,
  })

  const bot = botOptions.value.find(b => b.botID === data.botID)
  groupOptions.value = bot ? bot.botChatGroups : []

  if (formData.taskSendType === 3 || formData.taskSendType === 4) {
    let urls = []
    try { urls = JSON.parse(data.content || '[]') } catch(e){ urls=[] }
    formData.uploadFileList = urls.map(u => ({ name: '已上传文件', url: u }))
  }

  formData.extrendButton = data.extrendButton || []
  dialogFormVisible.value = true
}

// ================= 删除 ==================
const deleteBotTaskFunc = async (row) => {
  const res = await deleteBotTask({ ID: row.ID })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    if (tableData.value.length === 1 && page.value > 1) page.value--
    getTableData()
  }
}

// ================= 详情 ==================
const detailForm = ref({})
const detailShow = ref(false)
const openDetailShow = () => detailShow.value = true
const closeDetailShow = () => { detailShow.value = false; detailForm.value = {} }

const getDetails = async (row) => {
  const res = await findBotTask({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}

// ================= 显示发送类型文本 ==================
const renderSendType = (val) => {
  const n = Number(val)
  return { 1: '预设页面', 2: '文本', 3: '图片', 4: '视频' }[n] || ''
}

// ================= 解析 content JSON ==================
const parseContent = (content) => {
  try { return JSON.parse(content || '[]') } catch(e) { return [] }
}

// 扩展按钮二维数组初始化
if (!formData.extrendButton || !Array.isArray(formData.extrendButton)) {
  formData.extrendButton = [[]]
}

const dialogVisible = ref(false);
const editForm = ref({ name: "", url: "" });
let currentRow = null;
let currentIndex = null;

// 打开编辑弹窗
if (!formData.extrendButton || !Array.isArray(formData.extrendButton) || formData.extrendButton.length === 0) {
  formData.extrendButton = [[]];
}

// 打开编辑弹窗
function openEditDialog(rowIndex, btnIndex) {
  currentRow = rowIndex;
  currentIndex = btnIndex;
  editForm.value = { ...formData.extrendButton[rowIndex][btnIndex] };
  dialogVisible.value = true;
}

// 打开新增按钮弹窗（追加到当前行）
function openAddDialog(rowIndex) {
  currentRow = rowIndex;
  currentIndex = null; // 表示新增
  editForm.value = { name: "", url: "" };
  dialogVisible.value = true;
}

// 新增一行
function addNewRow() {
  formData.extrendButton.push([]); // 新增空数组作为新行
}

// 删除一行
function removeRow(rowIndex) {
  formData.extrendButton.splice(rowIndex, 1);
  if (formData.extrendButton.length === 0) {
    // 保证至少有一行
    formData.extrendButton.push([]);
  }
}

function removeButton(rowIndex, btnIndex) {
  formData.extrendButton[rowIndex].splice(btnIndex, 1);
}

// 保存按钮
function saveButton() {
  if (!editForm.value.name) return;
  if (currentIndex === null) {
    // 追加到当前行
    formData.extrendButton[currentRow].push({ ...editForm.value });
  } else {
    formData.extrendButton[currentRow][currentIndex] = { ...editForm.value };
  }
  dialogVisible.value = false;
}
</script>


<style scoped>
.btn-group-wrapper {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.btn-row {
  display: flex;
  flex-wrap: nowrap;       /* 强制不换行 */
  gap: 8px;                /* 按钮间距 */
  align-items: center;
  overflow-x: auto;        /* 超出宽度横向滚动 */
  padding: 5px 0;
}

.btn-item {
  display: flex;
  align-items: center;
  background-color: #f5f7fa;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 2px 8px;
  white-space: nowrap;
  gap: 5px;
  flex-shrink: 0;          /* 保证按钮不会缩小 */
}

.btn-text {
  font-size: 14px;
  line-height: 1;
}

.btn-edit, .btn-delete {
  cursor: pointer;
  color: #909399;
  transition: color 0.2s;
  font-size: 12px;
}

.btn-edit:hover {
  color: #409eff;
}

.btn-delete:hover {
  color: #f56c6c;
}

.add-btn {
  flex-shrink: 0;
  color: #409eff;
}

.delete-row-btn {
  flex-shrink: 0;
  color: #f56c6c;
}
</style>
