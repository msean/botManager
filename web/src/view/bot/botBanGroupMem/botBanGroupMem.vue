<template>
  <div>
    <!-- 搜索框 -->
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
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
            filterable
            clearable
            placeholder="请选择机器人"
            style="width: 220px"
          >
            <el-option
              v-for="item in botOptions"
              :key="item.value"
              :label="item.name"
              :value="item.botID"
            />
          </el-select>
        </el-form-item>
        

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格 -->
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">新增</el-button>
        <el-button
          icon="delete"
          style="margin-left: 10px;"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >
          删除
        </el-button>
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
        
        
        <el-table-column align="left" label="机器人名称" prop="botName" width="180" />
        <el-table-column align="left" label="群聊名称" prop="chatGroupName" width="360" />
        <el-table-column align="left" label="封禁成员" prop="banMemContent" width="180" />
        <el-table-column
          sortable
          align="left"
          label="日期"
          prop="createdAt"
          width="180"
        >
        <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="操作"
          fixed="right"
          :min-width="appStore.operateMinWith"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateBotBanGroupMemFunc(scope.row)"
            >
              编辑
            </el-button>
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
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

    <!-- 新增/编辑抽屉 -->
    <el-drawer
      destroy-on-close
      :size="appStore.drawerSize"
      v-model="dialogFormVisible"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增' : '编辑' }}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form
        :model="formData"
        label-position="top"
        ref="elFormRef"
        :rules="rule"
        label-width="80px"
      >
        <!-- 机器人选择 -->
        <el-form-item label="选择机器人" prop="botID">
          <el-select
            v-model="formData.botID"
            placeholder="请选择机器人"
            @change="handleBotChange"
            style="width: 100%"
            :disabled="isEdit"
          >
            <el-option
              v-for="bot in botOptions"
              :key="bot.botID"
              :label="bot.name"
              :value="bot.botID"
            />
          </el-select>
        </el-form-item>

        <!-- 群聊选择 -->
        <el-form-item label="选择群聊" prop="chatGroupID">
          <el-select
            v-model="formData.chatGroupID"
            placeholder="请选择群聊"
            :disabled="isEdit || !formData.botID"
            style="width: 100%"
          >
            <el-option
              v-for="chat in groupOptions"
              :key="chat.chatGroupID"
              :label="chat.chatGroupName"
              :value="chat.chatGroupID"
            />
          </el-select>
        </el-form-item>

        <!-- 封禁成员 -->
        <el-form-item label="封禁成员内容" prop="banMemContent">
          <el-input
            v-model="formData.banMemContent"
            :clearable="true"
            placeholder="请输入banMemContent"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createBotBanGroupMem,
  deleteBotBanGroupMem,
  deleteBotBanGroupMemByIds,
  updateBotBanGroupMem,
  findBotBanGroupMem,
  getBotBanGroupMemList
} from '@/api/bot/botBanGroupMem'
import { getBotChoiceWithChatGroup } from '@/api/bot/bot'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted, computed } from 'vue'
import { useAppStore } from "@/pinia"

defineOptions({ name: 'BotBanGroupMem' })

const btnLoading = ref(false)
const appStore = useAppStore()


// ========= 表格部分 =========
const formData = ref({ botID: undefined, chatGroupID: undefined, banMemContent: '' })
const rule = reactive({})
const elFormRef = ref()
const elSearchFormRef = ref()

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const multipleSelection = ref([])

const onReset = () => { searchInfo.value = {}; getTableData() }
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}
const handleSizeChange = (val) => { pageSize.value = val; getTableData() }
const handleCurrentChange = (val) => { page.value = val; getTableData() }

const getTableData = async() => {
  const table = await getBotBanGroupMemList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
  }
}
getTableData()

// ========== 新增/编辑 ==========
const dialogFormVisible = ref(false)
const type = ref('')
const isEdit = computed(() => type.value === 'update') // ✅ 是否编辑状态

const openDialog = () => {
  type.value = 'create'
  formData.value = { botID: undefined, chatGroupID: undefined, banMemContent: '' }
  groupOptions.value = []
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = { botID: undefined, chatGroupID: undefined, banMemContent: '' }
}

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async(valid) => {
    if (!valid) return (btnLoading.value = false)
    let res
    if (type.value === 'create') res = await createBotBanGroupMem(formData.value)
    else res = await updateBotBanGroupMem(formData.value)
    btnLoading.value = false
    if (res.code === 0) {
      ElMessage.success('保存成功')
      closeDialog()
      getTableData()
    }
  })
}

// ========== 删除与多选 ==========
const handleSelectionChange = (val) => (multipleSelection.value = val)
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { type: 'warning' }).then(() => deleteBotBanGroupMemFunc(row))
}
const deleteBotBanGroupMemFunc = async (row) => {
  const res = await deleteBotBanGroupMem({ ID: row.ID })
  if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
}
const onDelete = async() => {
  if (!multipleSelection.value.length) return ElMessage.warning('请选择要删除的数据')
  const IDs = multipleSelection.value.map(i => i.ID)
  const res = await deleteBotBanGroupMemByIds({ IDs })
  if (res.code === 0) { ElMessage.success('批量删除成功'); getTableData() }
}

// ========== 编辑 ==========
const updateBotBanGroupMemFunc = async (row) => {
  const res = await findBotBanGroupMem({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data
    type.value = 'update'
    // 设置群聊选项
    const selected = botOptions.value.find(b => b.botID === res.data.botID)
    groupOptions.value = selected ? selected.botChatGroups : []
    dialogFormVisible.value = true
  }
}

// ========== 联动下拉 ==========
const botOptions = ref([])
const groupOptions = ref([])

const loadBots = async () => {
  const res = await getBotChoiceWithChatGroup()
  if (res.code === 0) botOptions.value = res.data
}

const handleBotChange = (botID) => {
  const selected = botOptions.value.find(b => b.botID === botID)
  groupOptions.value = selected ? selected.botChatGroups : []
  formData.value.chatGroupID = undefined
}

onMounted(() => loadBots())
</script>

<style scoped>
.gva-table-box {
  margin-top: 20px;
}
</style>
