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
        <el-table-column prop="nextVerification" label="二步验证" />
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag v-if="scope.row.status === 0">未登录</el-tag>
            <el-tag type="warning" v-else-if="scope.row.status === 1">验证码</el-tag>
            <el-tag type="danger" v-else-if="scope.row.status === 2">二步验证</el-tag>
            <el-tag type="success" v-else>已登录</el-tag>
          </template>
        </el-table-column>

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
              disabled
            >自动二步验证</el-button>

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

        <el-form-item label="二步验证密码">
          <el-input
            v-model="formData.nextVerification"
            type="password"
            show-password
          />
        </el-form-item>

        <el-button type="primary" @click="save">保存</el-button>
      </el-form>
    </el-drawer>

    <!-- 验证码 -->
    <el-drawer v-model="codeVisible" title="输入验证码">

      <el-input v-model="loginCode" placeholder="Telegram 验证码" />

      <div style="margin-top:20px">
        <el-button type="primary" @click="submitCode">验证</el-button>

        <el-button
          type="warning"
          style="margin-left:10px"
          @click="resendCode"
        >
          重新发送验证码
        </el-button>
      </div>

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
  verifyCode
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
  apiHash: '',
  nextVerification: ''
})

const openDrawer = (row) => {

  formData.value = row
    ? { ...row }
    : {
        nickName: '',
        phone: '',
        apiId: 0,
        apiHash: '',
        nextVerification: ''
      }

  drawerVisible.value = true
}

const save = async () => {

  const api = formData.value.ID
    ? updateTgUser
    : createTgUser

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

const loginCode = ref('')


// 登录 → 发送验证码

const startLogin = async (row) => {

  try {

    const res = await sendCode({ id: row.ID })

    if (res.code === 0) {

      ElMessage.success(res.msg || '验证码已发送')

      currentRow.value = row

      codeVisible.value = true

      await getTableData()

    } else {

      ElMessage.error(res.msg || '发送验证码失败')

    }

  } catch (err) {

    ElMessage.error(err?.message || '发送验证码异常')

  }

}


// 打开验证码输入

const openCode = (row) => {

  currentRow.value = row

  codeVisible.value = true

}

// 重新发送验证码
const resendCode = async () => {

  if (!currentRow.value) return

  try {

    const res = await sendCode({
      id: currentRow.value.ID
    })

    if (res.code === 0) {

      ElMessage.success('验证码已重新发送')

    } else {

      ElMessage.error(res.msg || '发送失败')

    }

  } catch (err) {

    ElMessage.error('发送失败')

  }

}

// 验证验证码
const submitCode = async () => {

  const res = await verifyCode({
    id: currentRow.value.ID,
    code: loginCode.value
  })

  codeVisible.value = false
  loginCode.value = ''

  await getTableData()

  ElMessage.success('登录完成')
}

</script>