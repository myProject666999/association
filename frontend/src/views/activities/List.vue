<template>
  <div>
    <el-card class="common-card">
      <template #header>
        <div class="card-header">
          <span>活动列表</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            发布活动
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="活动名称" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部状态" clearable style="width: 150px">
            <el-option label="待审核" value="0" />
            <el-option label="已通过" value="1" />
            <el-option label="已拒绝" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="activityList" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="title" label="活动名称" min-width="200" />
        <el-table-column prop="club.name" label="所属社团" width="150">
          <template #default="scope">
            <el-tag v-if="scope.row.club" size="small" type="info">{{ scope.row.club.name }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="location" label="地点" width="120" />
        <el-table-column prop="start_time" label="开始时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.start_time) }}
          </template>
        </el-table-column>
        <el-table-column label="报名情况" width="120">
          <template #default="scope">
            {{ scope.row.current_participants }}/{{ scope.row.max_participants || '不限' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" link size="small" @click="handleViewDetail(scope.row)">
              详情
            </el-button>
            <el-button type="success" link size="small" @click="handleRegister(scope.row)" :disabled="scope.row.status !== 1">
              报名
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" title="发布活动" width="600px">
      <el-form ref="activityFormRef" :model="activityForm" :rules="activityRules" label-width="120px">
        <el-form-item label="活动名称" prop="title">
          <el-input v-model="activityForm.title" placeholder="请输入活动名称" />
        </el-form-item>
        <el-form-item label="活动地点" prop="location">
          <el-input v-model="activityForm.location" placeholder="请输入活动地点" />
        </el-form-item>
        <el-form-item label="活动描述" prop="description">
          <el-input
            v-model="activityForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入活动描述"
          />
        </el-form-item>
        <el-form-item label="活动开始时间" prop="start_time">
          <el-date-picker
            v-model="activityForm.start_time"
            type="datetime"
            placeholder="选择开始时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="活动结束时间" prop="end_time">
          <el-date-picker
            v-model="activityForm.end_time"
            type="datetime"
            placeholder="选择结束时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="报名开始时间" prop="registration_start">
          <el-date-picker
            v-model="activityForm.registration_start"
            type="datetime"
            placeholder="选择报名开始时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="报名结束时间" prop="registration_end">
          <el-date-picker
            v-model="activityForm.registration_end"
            type="datetime"
            placeholder="选择报名结束时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="最大参与人数" prop="max_participants">
          <el-input-number v-model="activityForm.max_participants" :min="0" placeholder="0表示不限" />
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">0表示不限人数</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/utils/request'

const router = useRouter()

const loading = ref(false)
const submitLoading = ref(false)
const activityList = ref([])
const dialogVisible = ref(false)

const searchForm = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const activityFormRef = ref(null)
const activityForm = reactive({
  title: '',
  description: '',
  location: '',
  start_time: null,
  end_time: null,
  registration_start: null,
  registration_end: null,
  max_participants: 0
})

const activityRules = {
  title: [{ required: true, message: '请输入活动名称', trigger: 'blur' }],
  location: [{ required: true, message: '请输入活动地点', trigger: 'blur' }],
  start_time: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  end_time: [{ required: true, message: '请选择结束时间', trigger: 'change' }],
  registration_start: [{ required: true, message: '请选择报名开始时间', trigger: 'change' }],
  registration_end: [{ required: true, message: '请选择报名结束时间', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.status) params.status = searchForm.status

    const res = await api.get('/api/activities', { params })
    if (res.code === 200) {
      activityList.value = res.data.list
      pagination.total = res.data.total
    }
  } catch (error) {
    console.error('获取活动列表失败', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  pagination.page = 1
  fetchData()
}

const handleCreate = () => {
  activityForm.title = ''
  activityForm.description = ''
  activityForm.location = ''
  activityForm.start_time = null
  activityForm.end_time = null
  activityForm.registration_start = null
  activityForm.registration_end = null
  activityForm.max_participants = 0
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!activityFormRef.value) return

  await activityFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        const data = { ...activityForm }
        const res = await api.post('/api/activities', data)
        if (res.code === 200) {
          ElMessage.success('活动发布成功，等待审核')
          dialogVisible.value = false
          fetchData()
        }
      } catch (error) {
        console.error('提交失败', error)
      } finally {
        submitLoading.value = false
      }
    }
  })
}

const handleViewDetail = (activity) => {
  router.push(`/activities/${activity.id}`)
}

const handleRegister = async (activity) => {
  try {
    ElMessageBox.confirm(`确定要报名参加"${activity.title}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }).then(async () => {
      const res = await api.post(`/api/activities/${activity.id}/register`)
      if (res.code === 200) {
        ElMessage.success('报名成功，等待审核')
      }
    }).catch(() => {})
  } catch (error) {
    console.error('报名失败', error)
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
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
