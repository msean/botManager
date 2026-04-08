<template>
  <div>
    <!-- 搜索 -->
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="机器人名称">
          <el-input v-model="searchInfo.name" />
        </el-form-item>

        <el-form-item label="机器人token">
          <el-input v-model="searchInfo.token" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格 -->
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" @click="openDialog">新增</el-button>
        <!-- <el-button :disabled="!multipleSelection.length" @click="onDelete">
          删除
        </el-button> -->
      </div>

      <el-table
        :data="tableData"
        row-key="botID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" prop="name" width="180" />
        <el-table-column label="Token" prop="token" width="420" />

        <!-- 是否记账 -->
        <el-table-column label="开启记账功能" width="140">
          <template #default="scope">
            <el-switch
              v-model="scope.row.isForLedger"
              :active-value="1"
              :inactive-value="2"
              @change="val => onLedgerSwitchChange(scope.row, val)"
            />
          </template>
        </el-table-column>

        <!-- ✅ 是否消息管理 -->
        <el-table-column label="开启消息管理功能" width="160">
          <template #default="scope">
            <el-switch
              v-model="scope.row.isForMsgMgr"
              :active-value="1"
              :inactive-value="2"
              @change="val => onMsgMgrSwitchChange(scope.row, val)"
            />
          </template>
        </el-table-column>

        <!-- ✅ 是否消息管理 -->
        <el-table-column label="开启群发消息" width="160">
          <template #default="scope">
            <el-switch
              v-model="scope.row.isForMsgMass"
              :active-value="1"
              :inactive-value="2"
              @change="val => onMsgMassSwitchChange(scope.row, val)"
            />
          </template>
        </el-table-column>

        <!-- 是否记账 -->
        <el-table-column label="自动广告发布" width="140">
          <template #default="scope">
            <el-switch
              v-model="scope.row.isAdPublish"
              :active-value="1"
              :inactive-value="2"
              @change="val => AdPublishSwitch(scope.row, val)"
            />
          </template>
        </el-table-column>

        <!-- <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column> -->

        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <!-- <el-button link @click="getDetails(scope.row)">查看</el-button> -->
            <el-button link @click="updateBotFunc(scope.row)">编辑</el-button>
            <el-button link type="danger" @click="deleteRow(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        layout="total, prev, pager, next"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 新增 / 编辑 -->
    <el-drawer v-model="dialogFormVisible" size="480px">
      <template #header>{{ type === 'create' ? '新增' : '编辑' }}</template>

      <el-form :model="formData" label-position="top">
        <el-form-item label="机器人名称">
          <el-input v-model="formData.name" />
        </el-form-item>

        <el-form-item label="Token">
          <el-input v-model="formData.token" />
        </el-form-item>

        <el-form-item label="是否记账">
          <el-switch
            v-model="formData.isForLedger"
            :active-value="1"
            :inactive-value="2"
          />
        </el-form-item>

        <el-form-item label="是否消息管理">
          <el-switch
            v-model="formData.isForMsgMgr"
            :active-value="1"
            :inactive-value="2"
          />
        </el-form-item>

        <el-form-item label="是否开启群发">
          <el-switch
            v-model="formData.isForMsgMass"
            :active-value="1"
            :inactive-value="2"
          />
        </el-form-item>

        <el-form-item label="是否开启群发">
          <el-switch
            v-model="formData.isForMsgMass"
            :active-value="1"
            :inactive-value="2"
          />
        </el-form-item>

         <el-form-item label="是否开启广告发布">
          <el-switch
            v-model="formData.isAdPublish"
            :active-value="1"
            :inactive-value="2"
          />
        </el-form-item>

        <el-button type="primary" @click="enterDialog">确定</el-button>
      </el-form>
    </el-drawer>

    <!-- 查看详情 -->
    <!-- <el-drawer v-model="detailShow" size="400px" title="查看机器人">
      <el-descriptions border :column="1">
        <el-descriptions-item label="名称">
          {{ detailForm.name }}
        </el-descriptions-item>
        <el-descriptions-item label="Token">
          {{ detailForm.token }}
        </el-descriptions-item>
        <el-descriptions-item label="是否记账">
          {{ detailForm.isForLedger === 1 ? '开启' : '关闭' }}
        </el-descriptions-item>
        <el-descriptions-item label="是否消息管理">
          {{ detailForm.isForMsgMgr === 1 ? '开启' : '关闭' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer> -->
  </div>
</template>

<script setup>
import {
  createBot,
  deleteBot,
  deleteBotByIds,
  updateBot,
  findBot,
  getBotList
} from '@/api/bot/bot'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'

defineOptions({ name: 'Bot' })

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const searchInfo = ref({})
const multipleSelection = ref([])

const dialogFormVisible = ref(false)
const detailShow = ref(false)
const type = ref('')

const formData = ref({
  name: '',
  token: '',
  isForLedger: 2,
  isForMsgMgr: 2
})

const detailForm = ref({})

const getTableData = async () => {
  const res = await getBotList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}
getTableData()

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const handleCurrentChange = (p) => {
  page.value = p
  getTableData()
}

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const openDialog = () => {
  type.value = 'create'
  formData.value = {
    name: '',
    token: '',
    isForLedger: 2,
    isForMsgMgr: 2
  }
  dialogFormVisible.value = true
}

const updateBotFunc = async (row) => {
  const res = await findBot({ ID: row.botID })
  if (res.code === 0) {
    formData.value = res.data
    type.value = 'update'
    dialogFormVisible.value = true
  }
}

const enterDialog = async () => {
  const api = type.value === 'create' ? createBot : updateBot
  const res = await api(formData.value)
  if (res.code === 0) {
    ElMessage.success('成功')
    dialogFormVisible.value = false
    getTableData()
  }
}

/** 是否记账 */
const onLedgerSwitchChange = async (row, newVal) => {
  const oldVal = newVal === 1 ? 2 : 1
  const res = await updateBot({
    botID: row.botID,
    isForLedger: newVal
  })
  if (res.code !== 0) {
    row.isForLedger = oldVal
    ElMessage.error('更新失败，已回滚')
  } else {
    ElMessage.success('更新成功')
  }
}


/** 开启广告自动发布 */
const AdPublishSwitch = async (row, newVal) => {
  const oldVal = newVal === 1 ? 2 : 1
  const res = await updateBot({
    botID: row.botID,
    isAdPublish: newVal
  })
  if (res.code !== 0) {
    row.isAdPublish = oldVal
    ElMessage.error('更新失败，已回滚')
  } else {
    ElMessage.success('更新成功')
  }
}

/** ✅ 是否消息管理 */
const onMsgMgrSwitchChange = async (row, newVal) => {
  const oldVal = newVal === 1 ? 2 : 1
  const res = await updateBot({
    botID: row.botID,
    isForMsgMgr: newVal
  })
  if (res.code !== 0) {
    row.isForMsgMgr = oldVal
    ElMessage.error('更新失败，已回滚')
  } else {
    ElMessage.success('更新成功')
  }
}


/** ✅ 是否消息管理 */
const onMsgMassSwitchChange = async (row, newVal) => {
  const oldVal = newVal === 1 ? 2 : 1
  const res = await updateBot({
    botID: row.botID,
    isForMsgMass: newVal
  })
  if (res.code !== 0) {
    row.isForMsgMass = oldVal
    ElMessage.error('更新失败，已回滚')
  } else {
    ElMessage.success('更新成功')
  }
}


const getDetails = async (row) => {
  const res = await findBot({ ID: row.botID })
  if (res.code === 0) {
    detailForm.value = res.data
    detailShow.value = true
  }
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定删除？').then(async () => {
    await deleteBot({ ID: row.botID })
    getTableData()
  })
}

const onDelete = async () => {
  const IDs = multipleSelection.value.map(i => i.botID)
  await deleteBotByIds({ IDs })
  getTableData()
}
</script>
