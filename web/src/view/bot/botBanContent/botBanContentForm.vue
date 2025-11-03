<template>
  <div>
    <div class="gva-form-box">
      <el-form
        :model="formData"
        ref="elFormRef"
        label-position="right"
        :rules="rule"
        label-width="80px"
      >
        <!-- 禁用内容 -->
        <el-form-item label="禁用内容:" prop="banContent">
          <el-input
            v-model="formData.banContent"
            clearable
            placeholder="请输入禁用内容"
          />
        </el-form-item>

        <!-- 机器人下拉选择 -->
        <el-form-item label="机器人:" prop="botID">
          <el-select
            v-model="formData.botID"
            placeholder="请选择机器"
            clearable
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="item in botList"
              :key="item.botID"
              :label="item.name"
              :value="item.botID"
            />
          </el-select>
        </el-form-item>

        <!-- 操作按钮 -->
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">
            保存
          </el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createBotBanContent,
  updateBotBanContent,
  findBotBanContent
} from '@/api/bot/botBanContent'

import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'
import request from '@/utils/request' // ✅ gin-vue-admin 内置的axios封装
import {
  getBotChoice,
} from '@/api/bot/bot'


defineOptions({
  name: 'botBanContentForm'
})

// ----------------- 数据区 -----------------
const route = useRoute()
const router = useRouter()
const btnLoading = ref(false)
const type = ref('')

const formData = ref({
  banContent: '',
  botID: undefined
})

const rule = reactive({
  banContent: [{ required: true, message: '请输入禁用内容', trigger: 'blur' }],
  botID: [{ required: true, message: '请选择机器人', trigger: 'change' }]
})

const elFormRef = ref()
const botList = ref([]) // 下拉列表数据

// ----------------- 方法区 -----------------

// 初始化：查询机器人列表 + 编辑时加载数据
const init = async () => {
  await getBotChoice()

  if (route.query.id) {
    const res = await findBotBanContent({ ID: route.query.id })
    if (res.code === 0) {
      formData.value = res.data
      type.value = 'update'
    }
  } else {
    type.value = 'create'
  }
}

// 保存按钮
const save = async () => {
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
      router.back()
    }
  })
}

// 返回
const back = () => router.go(-1)

onMounted(init)
</script>
