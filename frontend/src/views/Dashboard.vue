<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon" style="background: #409eff">
              <el-icon :size="30"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.users }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon" style="background: #67c23a">
              <el-icon :size="30"><OfficeBuilding /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.clubs }}</div>
              <div class="stat-label">社团总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon" style="background: #e6a23c">
              <el-icon :size="30"><Calendar /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.activities }}</div>
              <div class="stat-label">活动总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon" style="background: #f56c6c">
              <el-icon :size="30"><ChatDotRound /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.pending }}</div>
              <div class="stat-label">待审核</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card class="common-card">
          <template #header>
            <span>最新活动</span>
          </template>
          <el-table :data="latestActivities" stripe>
            <el-table-column prop="title" label="活动名称" min-width="200" />
            <el-table-column prop="location" label="地点" width="150" />
            <el-table-column label="时间" width="200">
              <template #default="scope">
                {{ formatDateTime(scope.row.start_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="current_participants" label="报名人数" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="common-card">
          <template #header>
            <span>快捷操作</span>
          </template>
          <el-row :gutter="10">
            <el-col :span="12" style="margin-bottom: 10px">
              <el-button type="primary" style="width: 100%" @click="goTo('/clubs')">
                <el-icon><Plus /></el-icon>
                浏览社团
              </el-button>
            </el-col>
            <el-col :span="12" style="margin-bottom: 10px">
              <el-button type="success" style="width: 100%" @click="goTo('/activities')">
                <el-icon><Plus /></el-icon>
                浏览活动
              </el-button>
            </el-col>
            <el-col :span="12" style="margin-bottom: 10px">
              <el-button type="warning" style="width: 100%" @click="goTo('/clubs/my')">
                <el-icon><OfficeBuilding /></el-icon>
                我的社团
              </el-button>
            </el-col>
            <el-col :span="12" style="margin-bottom: 10px">
              <el-button type="danger" style="width: 100%" @click="goTo('/activities/my')">
                <el-icon><Calendar /></el-icon>
                我的活动
              </el-button>
            </el-col>
          </el-row>
        </el-card>

        <el-card class="common-card">
          <template #header>
            <span>系统说明</span>
          </template>
          <div class="system-info">
            <p><strong>当前用户：</strong>{{ userStore.user?.name }}</p>
            <p><strong>角色：</strong>{{ getRoleText(userStore.user?.role) }}</p>
            <p><strong>院系：</strong>{{ userStore.user?.department || '未设置' }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import api from '@/utils/request'

const router = useRouter()
const userStore = useUserStore()

const stats = reactive({
  users: 0,
  clubs: 0,
  activities: 0,
  pending: 0
})

const latestActivities = ref([])

const fetchStats = async () => {
  try {
    const [usersRes, clubsRes, activitiesRes] = await Promise.all([
      api.get('/api/users?page_size=1'),
      api.get('/api/clubs?page_size=1'),
      api.get('/api/activities?page_size=1')
    ])
    stats.users = usersRes.data?.total || 0
    stats.clubs = clubsRes.data?.total || 0
    stats.activities = activitiesRes.data?.total || 0
  } catch (error) {
    console.error('获取统计数据失败', error)
  }
}

const fetchLatestActivities = async () => {
  try {
    const res = await api.get('/api/activities?page_size=5')
    latestActivities.value = res.data?.list || []
  } catch (error) {
    console.error('获取最新活动失败', error)
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
    case 2: return 'info'
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

const getRoleText = (role) => {
  switch (role) {
    case 'university_admin': return '校级管理员'
    case 'dept_admin': return '院级管理员'
    case 'student': return '学生'
    default: return '未知'
  }
}

const goTo = (path) => {
  router.push(path)
}

fetchStats()
fetchLatestActivities()
</script>

<style scoped>
.stat-card {
  margin-bottom: 20px;
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-info {
  margin-left: 15px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.system-info p {
  margin: 8px 0;
  color: #606266;
}
</style>
