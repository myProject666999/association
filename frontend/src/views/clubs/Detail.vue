<template>
  <div>
    <el-button type="primary" link @click="goBack" style="margin-bottom: 20px">
      <el-icon><ArrowLeft /></el-icon>
      返回列表
    </el-button>

    <el-card v-if="club" class="common-card">
      <template #header>
        <div class="card-header">
          <span>{{ club.name }}</span>
          <div>
            <el-tag :type="getStatusType(club.status)" size="small">
              {{ getStatusText(club.status) }}
            </el-tag>
            <el-tag v-if="club.category" size="small" type="info" style="margin-left: 10px">
              {{ club.category }}
            </el-tag>
          </div>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="社团名称">{{ club.name }}</el-descriptions-item>
        <el-descriptions-item label="分类">{{ club.category || '-' }}</el-descriptions-item>
        <el-descriptions-item label="成立时间">{{ club.founded_at ? formatDateTime(club.founded_at) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(club.status)">{{ getStatusText(club.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">
          {{ club.description || '暂无描述' }}
        </el-descriptions-item>
      </el-descriptions>

      <div style="margin-top: 20px">
        <el-button type="success" @click="handleApplyJoin" :disabled="club.status !== 1">
          申请加入
        </el-button>
      </div>
    </el-card>

    <el-card class="common-card" v-if="club">
      <template #header>
        <span>成员列表</span>
      </template>

      <el-table :data="clubMembers" stripe style="width: 100%" v-loading="membersLoading">
        <el-table-column prop="user.name" label="姓名" width="120">
          <template #default="scope">
            {{ scope.row.user?.name }}
          </template>
        </el-table-column>
        <el-table-column prop="user.student_id" label="学号" width="120">
          <template #default="scope">
            {{ scope.row.user?.student_id }}
          </template>
        </el-table-column>
        <el-table-column prop="position" label="职位" width="120">
          <template #default="scope">
            <el-tag :type="getPositionTagType(scope.row.position)" size="small">
              {{ getPositionText(scope.row.position) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
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
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/utils/request'

const route = useRoute()
const router = useRouter()

const club = ref(null)
const clubMembers = ref([])
const membersLoading = ref(false)

const clubId = route.params.id

const fetchClubDetail = async () => {
  try {
    const res = await api.get(`/api/clubs/${clubId}`)
    if (res.code === 200) {
      club.value = res.data
    }
  } catch (error) {
    console.error('获取社团详情失败', error)
  }
}

const fetchClubMembers = async () => {
  membersLoading.value = true
  try {
    const res = await api.get(`/api/clubs/${clubId}/members`)
    if (res.code === 200) {
      clubMembers.value = res.data || []
    }
  } catch (error) {
    console.error('获取社团成员失败', error)
  } finally {
    membersLoading.value = false
  }
}

const goBack = () => {
  router.back()
}

const handleApplyJoin = async () => {
  try {
    ElMessageBox.confirm(`确定要申请加入"${club.value.name}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }).then(async () => {
      const res = await api.post(`/api/clubs/${clubId}/apply`)
      if (res.code === 200) {
        ElMessage.success('申请已提交，请等待审核')
        fetchClubMembers()
      }
    }).catch(() => {})
  } catch (error) {
    console.error('申请加入失败', error)
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

onMounted(() => {
  fetchClubDetail()
  fetchClubMembers()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
