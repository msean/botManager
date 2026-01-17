<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const list = ref([])
const summary = ref({})
const loading = ref(true)

onMounted(async () => {
  const res = await axios.get('/api/ledger/full', {
    params: route.query
  })
  list.value = res.data.list
  summary.value = res.data.summary
  loading.value = false
})
</script>

<template>
  <div style="padding:16px">
    <h2>账单明细</h2>

    <el-table :data="list" v-loading="loading" border>
      <el-table-column label="时间" prop="created_at" width="160" />
      <el-table-column label="操作人" prop="opr_user_nickname" />
      <el-table-column label="金额" prop="amount" />
      <el-table-column label="输入" prop="raw_input" />
    </el-table>

    <el-divider />

    <el-descriptions border>
      <el-descriptions-item label="总入款">
        {{ summary.income }}
      </el-descriptions-item>
      <el-descriptions-item label="总下发">
        {{ summary.payout }}
      </el-descriptions-item>
      <el-descriptions-item label="未下发">
        {{ summary.unpaid }}
      </el-descriptions-item>
    </el-descriptions>
  </div>
</template>
