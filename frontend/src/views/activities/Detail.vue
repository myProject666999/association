<template>
  <div>
    <el-button type="primary" link @click="goBack" style="margin-bottom: 20px">
      <el-icon><ArrowLeft /></el-icon>
      返回列表
    </el-button>

    <el-card v-if="activity" class="common-card">
      <template #header>
        <div class="card-header">
          <span>{{ activity.title }}</span>
          <el-tag :type="getStatusType(activity.status)" size="small">
            {{ getStatusText(activity.status) }}
          </el-tag>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="活动名称">{{ activity.title }}</el-descriptions-item>
        <el-descriptions-item label="活动地点">{{ activity.location || '-' }}</el-descriptions-item>
        <el-descriptions-item label="开始时间">
          {{ formatDateTime(activity.start_time) }}
        </el-descriptions-item>
        <el-descriptions-item label="结束时间">
          {{ formatDateTime(activity.end_time) }}
        </el-descriptions-item>
        <el-descriptions-item label="报名开始时间">
          {{ formatDateTime(activity.registration_start) }}
        </el-descriptions-item>
        <el-descriptions-item label="报名结束时间">
          {{ formatDateTime(activity.registration_end) }}
        </el-descriptions-item>
        <el-descriptions-item label="参与人数">
          {{ activity.current_participants }} / {{ activity.max_participants || '不限' }}
        </el-descriptions-item>
        <el-descriptions-item label="所属社团">
          {{ activity.club?.name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="活动描述" :span="2">
          {{ activity.description || '暂无描述' }}
        </el-descriptions-item>
      </el-descriptions>

      <div style="margin-top: 20px">
        <el-button type="success" @click="handleRegister" :disabled="activity.status !== 1">
          报名参加
        </el-button>
      </div>
    </el-card>

    <el-card class="common-card" v-if="activity">
      <template #header>
        <span>报名列表</span>
      </template>

      <el-table :data="registrations" stripe style="width: 100%" v-loading="registrationsLoading">
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
        <el-table-column prop="status" label="报名状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="registrated_at" label="报名时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.registrated_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="common-card" v-if="activity">
      <template #header>
        <span>活动评论</span>
      </template>

      <el-form v-if="activity.status === 1" style="margin-bottom: 20px">
        <el-input
          v-model="commentForm.content"
          type="textarea"
          :rows="2"
          placeholder="发表评论..."
          maxlength="500"
          show-word-limit
        />
        <div style="margin-top: 10px; text-align: right">
          <el-button type="primary" :loading="commentLoading" @click="handleAddComment">
            发表评论
          </el-button>
        </div>
      </el-form>

      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <div class="comment-header">
          <el-avatar :size="32">
          {{ comment.user?.name?.charAt(0) }}
          </el-avatar>
          <div class="comment-info">
            <span class="comment-author">{{ comment.user?.name }}</span>
            <span class="comment-time">{{ formatDateTime(comment.created_at) }}</span>
          </div>
        </div>
        <div class="comment-content">{{ comment.content }}</div>
      </div>

      <el-empty v-if="comments.length === 0" description="暂无评论" />
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

const activity = ref(null)
const registrations = ref([])
const comments = ref([])
const registrationsLoading = ref(false)
const commentLoading = ref(false)

const commentForm = reactive({
  content: ''
})

const activityId = route.params.id

const fetchActivityDetail = async () => {
  try {
    const res = await api.get(`/api/activities/${activityId}`)
    if (res.code === 200) {
      activity.value = res.data
    }
  } catch (error) {
    console.error('获取活动详情失败', error)
  }
}

const fetchRegistrations = async () => {
  registrationsLoading.value = true
  try {
    const res = await api.get(`/api/activities/${activityId}/registrations`)
    if (res.code === 200) {
      registrations.value = res.data || []
    }
  } catch (error) {
    console.error('获取报名列表失败', error)
  } finally {
    registrationsLoading.value = false
  }
}

const fetchComments = async () => {
  try {
    const res = await api.get(`/api/activities/${activityId}/comments`)
    if (res.code === 200) {
      comments.value = res.data || []
    }
  } catch (error) {
    console.error('获取评论列表失败', error)
  }
}

const goBack = () => {
  router.back()
}

const handleRegister = async () => {
  try {
    ElMessageBox.confirm(`确定要报名参加"${activity.value.title}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }).then(async () => {
      const res = await api.post(`/api/activities/${activityId}/register`)
      if (res.code === 200) {
        ElMessage.success('报名成功，等待审核')
        fetchRegistrations()
      }
    }).catch(() => {})
  } catch (error) {
    console.error('报名失败', error)
  }
}

const handleAddComment = async () => {
  if (!commentForm.content.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  commentLoading.value = true
  try {
    const res = await api.post(`/api/activities/${activityId}/comments`, {
      content: commentForm.content
    })
    if (res.code === 200) {
      ElMessage.success('评论成功')
      commentForm.content = ''
      fetchComments()
    }
  } catch (error) {
    console.error('评论失败', error)
  } finally {
    commentLoading.value = false
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
  fetchActivityDetail()
  fetchRegistrations()
  fetchComments()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.comment-item {
  padding: 15px 0;
  border-bottom: 1px solid #ebeef5;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  align-items: center;
}

.comment-info {
  margin-left: 10px;
}

.comment-author {
  font-weight: 500;
  margin-right: 10px;
}

.comment-time {
  color: #909399;
  font-size: 12px;
}

.comment-content {
  margin-top: 8px;
  color: #606266;
  line-height: 1.6;
}
</style>
