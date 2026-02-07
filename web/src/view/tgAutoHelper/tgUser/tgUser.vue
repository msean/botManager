<template>
  <div>
    <!-- 搜索 -->
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="昵称">
          <el-input v-model="searchInfo.nickName" placeholder="昵称" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getTableData">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格 -->
    <div class="gva-table-box">
      <el-button type="primary" @click="openDrawer()">新增</el-button>

      <el-table :data="tableData" row-key="ID">
        <el-table-column prop="nickName" label="昵称" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="apiId" label="ApiId" />
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag v-if="scope.row.status === 0">未登录</el-tag>
            <el-tag type="warning" v-else-if="scope.row.status === 1">验证码</el-tag>
            <el-tag type="danger" v-else-if="scope.row.status === 2">二步验证</el-tag>
            <el-tag type="success" v-else>已登录</el-tag>
          </template>
        </el-table-column>

        <!-- 操作 -->
        <el-table-column fixed="right" label="操作" width="260">
          <template #default="scope">
            <el-button
              v-if="scope.row.status === 0"
              link
              type="primary"
              @click="startLogin(scope.row)"
            >登录</el-button>

            <el-button
              v-if="scope.row.status === 1"
              link
              type="warning"
              @click="openCode(scope.row)"
            >验证码</el-button>

            <el-button
              v-if="scope.row.status === 2"
              link
              type="danger"
              @click="openPassword(scope.row)"
            >二步验证</el-button>

            <el-tag v-if="scope.row.status === 3" type="success">
              ✓
            </el-tag>

            <el-divider direction="vertical" />

            <el-button link @click="openDrawer(scope.row)">编辑</el-button>
            <el-button link type="danger" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        layout="total, sizes, prev, pager, next"
        :total="total"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        @current-change="getTableData"
        @size-change="getTableData"
      />
    </div>

    <!-- 新增 / 编辑 -->
    <el-drawer v-model="drawerVisible" title="Telegram 用户">
      <el-form :model="formData" label-position="top">
        <el-form-item label="昵称">
          <el-input v-model="formData.nickName" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="formData.phone" />
        </el-form-item>
        <el-form-item label="ApiId">
          <el-input v-model.number="formData.apiId" />
        </el-form-item>
        <el-form-item label="ApiHash">
          <el-input v-model="formData.apiHash" />
        </el-form-item>
        <el-button type="primary" @click="save">保存</el-button>
      </el-form>
    </el-drawer>

    <!-- 验证码 -->
    <el-drawer v-model="codeVisible" title="输入验证码">
      <el-input v-model="loginCode" placeholder="Telegram 验证码" />
      <el-button type="primary" @click="submitCode">验证</el-button>
    </el-drawer>

    <!-- 二步验证 -->
    <el-drawer v-model="passwordVisible" title="二步验证">
      <el-input
        v-model="loginPassword"
        type="password"
        placeholder="二步验证密码"
      />
      <el-button type="primary" @click="submitPassword">验证</el-button>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getTgUserList,
  createTgUser,
  updateTgUser,
  deleteTgUser,
  sendCode,
  verifyCode,
  verifyPassword
} from '@/api/tg_auto_helper/tgUser'

// ===== 列表 =====
const searchInfo = ref({})
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const getTableData = async () => {
  const res = await getTgUserList({
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

const resetSearch = () => {
  searchInfo.value = {}
  getTableData()
}

// ===== 新增 / 编辑 =====
const drawerVisible = ref(false)
const formData = ref({
  nickName: '',
  phone: '',
  apiId: 0,
  apiHash: ''
})

const openDrawer = (row) => {
  formData.value = row ? { ...row } : {
    nickName: '',
    phone: '',
    apiId: 0,
    apiHash: ''
  }
  drawerVisible.value = true
}

const save = async () => {
  const api = formData.value.ID ? updateTgUser : createTgUser
  const res = await api(formData.value)
  if (res.code === 0) {
    ElMessage.success('成功')
    drawerVisible.value = false
    getTableData()
  }
}

const deleteRow = async (row) => {
  await ElMessageBox.confirm('确认删除？')
  await deleteTgUser({ ID: row.ID })
  ElMessage.success('已删除')
  getTableData()
}

// ===== 登录流程 =====
const currentRow = ref(null)
const codeVisible = ref(false)
const passwordVisible = ref(false)
const loginCode = ref('')
const loginPassword = ref('')

// 登录 → 发送验证码
const startLogin = async (row) => {
  try {
    const res = await sendCode({ id: row.ID })

    // ✅ 只有成功才继续
    if (res.code === 0) {
      ElMessage.success(res.msg || '验证码已发送')

      currentRow.value = row
      codeVisible.value = true
      await getTableData()
    } else {
      ElMessage.error(res.msg || '发送验证码失败')
    }
  } catch (err) {
    // ✅ 网络 / 500 / 超时
    ElMessage.error(err?.message || '发送验证码异常')
  }
}


const openCode = (row) => {
  currentRow.value = row
  codeVisible.value = true
}

// 验证验证码
const submitCode = async () => {
  const res = await verifyCode({
    id: currentRow.value.ID,
    code: loginCode.value
  })

  ElMessage.success('验证码正确')
  codeVisible.value = false
  loginCode.value = ''

  // ⭐ 重新拉数据，看是否进入二步
  await getTableData()

  const updated = tableData.value.find(
    (u) => u.ID === currentRow.value.ID
  )

  if (updated && updated.status === 2) {
    passwordVisible.value = true
  }
}

// 二步验证
const openPassword = (row) => {
  currentRow.value = row
  passwordVisible.value = true
}

const submitPassword = async () => {
  await verifyPassword({
    id: currentRow.value.ID,
    password: loginPassword.value
  })

  ElMessage.success('登录完成')
  passwordVisible.value = false
  loginPassword.value = ''
  getTableData()
}
</script>
