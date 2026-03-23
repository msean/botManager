<template>
  <div>

    <!-- 搜索 -->
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">

        <el-form-item label="关键词">
          <el-input v-model="searchInfo.searchText" placeholder="关键词"/>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="getTableData">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>

      </el-form>
    </div>


    <!-- 表格 -->
    <div class="gva-table-box">

      <el-button type="primary" @click="openDrawer()">新增任务</el-button>

      <el-table :data="tableData" row-key="ID">

        <el-table-column prop="ID" label="ID" width="80"/>

        <el-table-column prop="searchText" label="关键词"/>

        <el-table-column prop="memberLower" label="人数下限"/>

        <el-table-column prop="totalPage" label="总页数"/>

        <el-table-column prop="currentPage" label="当前页"/>

        <el-table-column prop="totalCount" label="采集群数"/>

        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag v-if="scope.row.status === 1" type="warning">
              采集中
            </el-tag>
            <el-tag v-else-if="scope.row.status === 2" type="success">
              完成
            </el-tag>
            <el-tag v-else type="danger">
              停止
            </el-tag>
          </template>
        </el-table-column>


        <!-- ✅ 操作列（已新增查看按钮） -->
        <el-table-column fixed="right" label="操作" width="220">
          <template #default="scope">

            <el-button
              link
              type="primary"
              @click="openDetail(scope.row)"
            >
              查看数据
            </el-button>

            <el-button
              link
              type="danger"
              @click="deleteRow(scope.row)"
            >
              删除
            </el-button>

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


    <!-- 新增任务 -->
    <el-drawer v-model="drawerVisible" title="创建采集任务">
      <el-form :model="formData" label-position="top">

        <el-form-item label="采集关键词">
          <el-input
            type="textarea"
            :rows="1"
            v-model="formData.searchText"
          />
        </el-form-item>

        <el-form-item label="群人数下限(0 表示无成员数量限制)">
          <el-input-number v-model="formData.memberLower" :min="0"/>
        </el-form-item>

        <el-button type="primary" @click="save">
          创建任务
        </el-button>

      </el-form>
    </el-drawer>


    <!-- ✅ 采集数据弹窗 -->
    <el-drawer v-model="detailVisible" title="采集数据" size="60%">

      <!-- 搜索 -->
      <div style="margin-bottom: 10px">
        <el-input
          v-model="detailSearchText"
          placeholder="搜索群名称"
          style="width: 200px"
        />
        <el-button type="primary" @click="getDetailList">搜索</el-button>
      </div>

      <!-- 表格 -->
      <el-table :data="detailList">

        <el-table-column prop="ID" label="ID" width="80"/>

        <el-table-column prop="title" label="群名称"/>

        <el-table-column prop="link" label="链接">
          <template #default="scope">
            <el-link :href="scope.row.link" target="_blank">
              打开
            </el-link>
          </template>
        </el-table-column>

        <el-table-column prop="members" label="人数" sortable/>

      </el-table>

      <!-- 分页 -->
      <el-pagination
        layout="total, prev, pager, next"
        :total="detailTotal"
        v-model:current-page="detailPage"
        v-model:page-size="detailPageSize"
        @current-change="getDetailList"
      />

    </el-drawer>

  </div>
</template>


<script setup>

import {ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'

import {
  createCollectGroupTask,
  deleteCollectGroupTask,
  listCollectGroupTask,
  listCollectGroupInfo
} from '@/api/tg_auto_helper/collectGroup'


// ================= 列表 =================

const searchInfo = ref({})
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const getTableData = async () => {

  const res = await listCollectGroupTask({
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


// ================= 新增 =================

const drawerVisible = ref(false)

const formData = ref({
  searchText: '',
  memberLower: 0
})

const openDrawer = () => {
  formData.value = {
    searchText: '',
    memberLower: 0
  }
  drawerVisible.value = true
}

const save = async () => {
  const res = await createCollectGroupTask(formData.value)
  if (res.code === 0) {
    ElMessage.success('创建成功')
    drawerVisible.value = false
    getTableData()
  }
}


// ================= 删除 =================

const deleteRow = async (row) => {
  await ElMessageBox.confirm('确认删除任务？')
  await deleteCollectGroupTask({ ID: row.ID })
  ElMessage.success('删除成功')
  getTableData()
}


// ================= ✅ 采集数据 =================

const detailVisible = ref(false)
const detailList = ref([])
const detailTotal = ref(0)
const detailPage = ref(1)
const detailPageSize = ref(10)
const currentTaskID = ref(0)
const detailSearchText = ref('')

// 打开弹窗
const openDetail = (row) => {
  currentTaskID.value = row.ID
  detailPage.value = 1
  detailVisible.value = true
  getDetailList()
}

// 获取数据
const getDetailList = async () => {

  const res = await listCollectGroupInfo({
    page: detailPage.value,
    pageSize: detailPageSize.value,
    taskID: currentTaskID.value,
    searchText: detailSearchText.value
  })

  if (res.code === 0) {
    detailList.value = res.data.list
    detailTotal.value = res.data.total
  }
}

</script>