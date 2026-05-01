<template>
  <div>
    <el-card class="common-card">
      <template #header>
        <span>我的活动</span>
      </template>

      <el-table :data="activityList" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="activity.title" label="活动名称" min-width="200">
          <template #default="scope">
            {{ scope.row.activity?.title }}
          </template>
        </el-table-column>
        <el-table-column prop="activity.club.name" label="所属社团" width="150">
          <template #default="scope">
            <el-tag v-if="scope.row.activity?.club" size="small" type="info">
              {{ scope.row.activity.club.name }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="activity.location" label="地点" width="120">
          <template #default="scope">
            {{ scope.row.activity?.location }}
          </template>
        </el-table-column>
        <el-table-column label="开始时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.activity?.start_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="报名状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="报名时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button type="primary" link size="small" @click="handleViewDetail(scope.row.activity)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && activityList.length === 0" description="暂无报名的活动" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/request'

const router = useRouter()

const loading = ref(false)
const activityList = ref([])

const fetchData = async () => {
  loading.value = true
  try {
    const res = await api.get('/api/profile/activities')
    if (res.code === 200) {
      activityList.value = res.data || []
    }
  } catch (error) {
    console.error('获取我的活动失败', error)
  } finally {
    loading.value = false
  }
}

const handleViewDetail = (activity) => {
  if (activity) {
    router.push(`/activities/${activity.id}`)
  }
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const getStatusType = (status) => {
  switch (status) {
    case 0: return 'warning'
    case 1: return 'success'
    case 2: return 'danger'
    default: return 'info'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 0: return '待审核'
    case 1: return '已通过'
    case 2: return '已拒绝'
    default: return '未知'
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
</style>
