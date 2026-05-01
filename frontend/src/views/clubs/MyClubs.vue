<template>
  <div>
    <el-card class="common-card">
      <template #header>
        <span>我的社团</span>
      </template>

      <el-table :data="clubList" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="club.name" label="社团名称" width="200">
          <template #default="scope">
            {{ scope.row.club?.name }}
          </template>
        </el-table-column>
        <el-table-column prop="club.category" label="分类" width="120">
          <template #default="scope">
            <el-tag size="small">{{ scope.row.club?.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="position" label="职位" width="120">
          <template #default="scope">
            <el-tag :type="getPositionTagType(scope.row.position)" size="small">
              {{ getPositionText(scope.row.position) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="joined_at" label="加入时间" width="180">
          <template #default="scope">
            {{ scope.row.joined_at ? formatDateTime(scope.row.joined_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button type="primary" link size="small" @click="handleViewDetail(scope.row.club)">
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/request'

const router = useRouter()

const loading = ref(false)
const clubList = ref([])

const fetchData = async () => {
  loading.value = true
  try {
    const res = await api.get('/api/profile/clubs')
    if (res.code === 200) {
      clubList.value = res.data || []
    }
  } catch (error) {
    console.error('获取我的社团失败', error)
  } finally {
    loading.value = false
  }
}

const handleViewDetail = (club) => {
  if (club) {
    router.push(`/clubs/${club.id}`)
  }
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const getPositionTagType = (position) => {
  switch (position) {
    case 'president': return 'danger'
    case 'vice_president': return 'warning'
    default: return 'info'
  }
}

const getPositionText = (position) => {
  switch (position) {
    case 'president': return '社长'
    case 'vice_president': return '副社长'
    default: return '成员'
  }
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
