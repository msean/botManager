<template>
  <div class="gva-container">
    <!-- 表格 -->
    <el-table :data="tableData" stripe>
      <el-table-column prop="nickName" label="昵称" />
      <el-table-column prop="phone" label="手机号" />
      <el-table-column prop="apiId" label="ApiId" />
      <el-table-column prop="status" label="状态">
        <template #default="{ row }">
          <el-tag v-if="row.status === 0">未登录</el-tag>
          <el-tag type="warning" v-else-if="row.status === 1">已发送验证码</el-tag>
          <el-tag type="info" v-else-if="row.status === 2">二步验证</el-tag>
          <el-tag type="success" v-else>已登录</el-tag>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="260">
        <template #default="{ row }">
          <el-button size="small" @click="openDrawer(row)">编辑</el-button>
          <el-button
            size="small"
            type="primary"
            v-if="row.status === 0"
            @click="handleSendCode(row)"
          >
            登录
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 抽屉 -->
    <el-drawer v-model="drawerVisible" title="Telegram 登录" size="40%">
      <el-form :model="formData" label-width="100px">
        <el-form-item label="昵称">
          <el-input v-model="formData.nickName" />
        </el-form-item>

        <el-form-item label="手机号">
          <el-input v-model="formData.phone" />
        </el-form-item>

        <el-form-item label="ApiId">
          <el-input-number
            v-model="formData.apiId"
            :min="0"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="ApiHash">
          <el-input v-model="formData.apiHash" />
        </el-form-item>

        <!-- 验证码 -->
        <el-form-item label="验证码" v-if="formData.status === 1">
          <el-input v-model="verifyCodeValue" />
        </el-form-item>

        <!-- 二步密码 -->
        <el-form-item label="二步密码" v-if="formData.status === 2">
          <el-input v-model="passwordValue" type="password" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>

        <el-button type="primary" @click="handleSave">
          保存
        </el-button>

        <el-button
          type="primary"
          v-if="formData.status === 0"
          @click="handleSendCode(formData)"
        >
          发送验证码
        </el-button>

        <el-button
          type="success"
          v-if="formData.status === 1"
          @click="handleVerifyCode"
        >
          验证验证码
        </el-button>

        <el-button
          type="success"
          v-if="formData.status === 2"
          @click="handleVerifyPassword"
        >
          二步验证
        </el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  getTgUserList,
  saveTgUser,
  sendCode,
  verifyCode,
  verifyPassword
} from '@/api/tgUser'
import { ElMessage } from 'element-plus'

/* ===== 数据 ===== */
const tableData = ref([])
const drawerVisible = ref(false)

const emptyForm = {
  ID: '',
  nickName: '',
  phone: '',
  apiId: 0,
  apiHash: '',
  status: 0
}

const formData = ref({ ...emptyForm })

const verifyCodeValue = ref('')
const passwordValue = ref('')

/* ===== 方法 ===== */
const loadData = async () => {
  const res = await getTgUserList()
  if (res.code === 0) {
    tableData.value = res.data.list
  }
}

const openDrawer = (row) => {
  formData.value = row
    ? { ...emptyForm, ...row }
    : { ...emptyForm }

  verifyCodeValue.value = ''
  passwordValue.value = ''
  drawerVisible.value = true
}

const handleSave = async () => {
  const res = await saveTgUser(formData.value)
  if (res.code === 0) {
    ElMessage.success('保存成功')
    drawerVisible.value = false
    loadData()
  }
}

const handleSendCode = async (row) => {
  const res = await sendCode({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success('验证码已发送')
    row.status = 1
  }
}

const handleVerifyCode = async () => {
  const res = await verifyCode({
    id: formData.value.ID,
    code: verifyCodeValue.value
  })
  if (res.code === 0) {
    ElMessage.success('验证码通过')
    formData.value.status = res.data.status
    loadData()
  }
}

const handleVerifyPassword = async () => {
  const res = await verifyPassword({
    id: formData.value.ID,
    password: passwordValue.value
  })
  if (res.code === 0) {
    ElMessage.success('登录成功')
    formData.value.status = 3
    drawerVisible.value = false
    loadData()
  }
}

onMounted(loadData)
</script>

<style scoped>
.gva-container {
  padding: 16px;
}
</style>
