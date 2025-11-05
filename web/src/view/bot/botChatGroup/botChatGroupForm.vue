
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <!-- 机器人下拉选择 -->
        <el-form-item label="机器人:" prop="botID">
          <el-select
            v-model="formData.botID"
            placeholder="请选择机器人"
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
        <el-form-item label="群组ID:" prop="chatGroupID">
            <el-input v-model.number="formData.chatGroupID" :clearable="true" placeholder="请输入群组ID" />
        </el-form-item>
        <el-form-item label="群组名称:" prop="chatGroupName">
                <el-input v-model.number="formData.chatGroupName" :clearable="true" placeholder="请输入群组名称" />
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
  createBotChatGroup,
  updateBotChatGroup,
  findBotChatGroup
} from '@/api/bot/botChatGroup'

defineOptions({
    name: 'BotChatGroupForm'
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
            botID: undefined,
            chatGroupID: undefined,
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findBotChatGroup({ ID: route.query.id })
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
               res = await createBotChatGroup(formData.value)
               break
             case 'update':
               res = await updateBotChatGroup(formData.value)
               break
             default:
               res = await createBotChatGroup(formData.value)
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
