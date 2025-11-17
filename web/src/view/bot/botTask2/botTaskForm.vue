
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="机器人:" prop="botID">
    <el-input v-model.number="formData.botID" :clearable="true" placeholder="请输入机器人ID" />
</el-form-item>
        <el-form-item label="群聊:" prop="chatGroupID">
    <el-input v-model.number="formData.chatGroupID" :clearable="true" placeholder="请输入群ID" />
</el-form-item>
        <el-form-item label="发送类型:" prop="taskSendType">
    <el-input v-model.number="formData.taskSendType" :clearable="true" placeholder="请输入发送类型" />
</el-form-item>
        <el-form-item label="发送内容:" prop="content">
    <RichEdit v-model="formData.content"/>
</el-form-item>
        <el-form-item label="扩展按钮:" prop="extrendButton">
    <RichEdit v-model="formData.extrendButton"/>
</el-form-item>
        <el-form-item label="发送间隔:" prop="sendInterval">
    <el-input v-model.number="formData.sendInterval" :clearable="true" placeholder="请输入发送间隔" />
</el-form-item>
        <el-form-item label="下一次发送时间:" prop="nextSendTime">
    <el-date-picker v-model="formData.nextSendTime" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
        <el-form-item label="上一次发送时间:" prop="preSendTime">
    <el-date-picker v-model="formData.preSendTime" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
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
  createBotTask,
  updateBotTask,
  findBotTask
} from '@/api/bot/botTask'

defineOptions({
    name: 'BotTaskForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const formData = ref({
            botID: undefined,
            chatGroupID: undefined,
            taskSendType: undefined,
            content: '',
            extrendButton: '',
            sendInterval: undefined,
            nextSendTime: new Date(),
            preSendTime: new Date(),
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
      const res = await findBotTask({ ID: route.query.id })
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
               res = await createBotTask(formData.value)
               break
             case 'update':
               res = await updateBotTask(formData.value)
               break
             default:
               res = await createBotTask(formData.value)
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
