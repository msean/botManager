
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="昵称:" prop="nickName">
    <el-input v-model="formData.nickName" :clearable="true" placeholder="请输入昵称" />
</el-form-item>
        <el-form-item label="手机号码:" prop="phone">
    <el-input v-model="formData.phone" :clearable="true" placeholder="请输入手机号码" />
</el-form-item>
        <el-form-item label="apiId:" prop="apiId">
    <el-input v-model="formData.apiId" :clearable="true" placeholder="请输入apiId" />
</el-form-item>
        <el-form-item label="apiHash:" prop="apiHash">
    <el-input v-model="formData.apiHash" :clearable="true" placeholder="请输入apiHash" />
</el-form-item>
        <el-form-item label="二步验证:" prop="nextVerification">
    <el-input v-model="formData.nextVerification" :clearable="true" placeholder="请输入二步验证" />
</el-form-item>
        <el-form-item label="状态:" prop="status">
    <el-input v-model.number="formData.status" :clearable="true" placeholder="请输入状态" />
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createTgUser,
  updateTgUser,
  findTgUser
} from '@/api/tg_auto_helper/tgUser'

defineOptions({
    name: 'TgUserForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const formData = ref({
            nickName: '',
            phone: '',
            apiId: 0,
            apiHash: '',
            nextVerification: '',
            status: undefined,
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findTgUser({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createTgUser(formData.value)
               break
             case 'update':
               res = await updateTgUser(formData.value)
               break
             default:
               res = await createTgUser(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
