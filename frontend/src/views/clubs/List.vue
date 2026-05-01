<template>
  <div>
    <el-card class="common-card">
      <template #header>
        <div class="card-header">
          <span>社团列表</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            申请创建社团
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="社团名称" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.category" placeholder="全部分类" clearable style="width: 150px">
            <el-option label="学术科技" value="学术科技" />
            <el-option label="文化艺术" value="文化艺术" />
            <el-option label="体育竞技" value="体育竞技" />
            <el-option label="志愿服务" value="志愿服务" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-row :gutter="20">
        <el-col :span="8" v-for="club in clubList" :key="club.id">
          <el-card class="club-card" shadow="hover">
            <div class="club-header">
              <el-avatar :size="60">
                <el-icon v-if="!club.logo" :size="40"><OfficeBuilding /></el-icon>
                <img v-else :src="club.logo" />
              </el-avatar>
              <div class="club-info">
                <h3>{{ club.name }}</h3>
                <el-tag :type="getStatusType(club.status)" size="small">
                  {{ getStatusText(club.status) }}
                </el-tag>
              </div>
            </div>
            <div class="club-body">
              <p class="club-desc">{{ club.description || '暂无描述' }}</p>
              <div class="club-footer">
                <el-tag size="small" v-if="club.category">{{ club.category }}</el-tag>
                <span>成员: {{ club.members?.length || 0 }}</span>
              </div>
            </div>
            <div class="club-actions">
              <el-button type="primary" size="small" @click="handleViewDetail(club)">查看详情</el-button>
              <el-button type="success" size="small" @click="handleApplyJoin(club)">申请加入</el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[12, 24, 48]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" title="创建社团申请" width="500px">
      <el-form ref="clubFormRef" :model="clubForm" :rules="clubRules" label-width="80px">
        <el-form-item label="社团名称" prop="name">
          <el-input v-model="clubForm.name" placeholder="请输入社团名称" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="clubForm.category" placeholder="请选择分类" style="width: 100%">
            <el-option label="学术科技" value="学术科技" />
            <el-option label="文化艺术" value="文化艺术" />
            <el-option label="体育竞技" value="体育竞技" />
            <el-option label="志愿服务" value="志愿服务" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="clubForm.description"
            type="textarea"
            :rows="4"
            placeholder="请输入社团描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">提交申请</el-button>
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
const clubList = ref([])
const dialogVisible = ref(false)

const searchForm = reactive({
  keyword: '',
  category: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 12,
  total: 0
})

const clubFormRef = ref(null)
const clubForm = reactive({
  name: '',
  category: '',
  description: ''
})

const clubRules = {
  name: [{ required: true, message: '请输入社团名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.category) params.category = searchForm.category

    const res = await api.get('/api/clubs', { params })
    if (res.code === 200) {
      clubList.value = res.data.list
      pagination.total = res.data.total
    }
  } catch (error) {
    console.error('获取社团列表失败', error)
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
  searchForm.category = ''
  pagination.page = 1
  fetchData()
}

const handleCreate = () => {
  clubForm.name = ''
  clubForm.category = ''
  clubForm.description = ''
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!clubFormRef.value) return

  await clubFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        const res = await api.post('/api/clubs', clubForm)
        if (res.code === 200) {
          ElMessage.success('社团创建申请已提交')
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

const handleViewDetail = (club) => {
  router.push(`/clubs/${club.id}`)
}

const handleApplyJoin = async (club) => {
  try {
    ElMessageBox.confirm(`确定要申请加入"${club.name}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }).then(async () => {
      const res = await api.post(`/api/clubs/${club.id}/apply`)
      if (res.code === 200) {
        ElMessage.success('申请已提交，请等待审核')
      }
    }).catch(() => {})
  } catch (error) {
    console.error('申请加入失败', error)
  }
}

const getStatusType = (status) => {
  switch (status) {
    case 0: return 'warning'
    case 1: return 'success'
    default: return 'info'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 0: return '待审核'
    case 1: return '已启用'
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

.club-card {
  margin-bottom: 20px;
}

.club-header {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.club-info {
  margin-left: 15px;
}

.club-info h3 {
  margin: 0 0 5px 0;
  font-size: 16px;
  color: #303133;
}

.club-body {
  margin-bottom: 15px;
}

.club-desc {
  margin: 0;
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.club-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
  font-size: 12px;
  color: #909399;
}

.club-actions {
  display: flex;
  gap: 10px;
  padding-top: 15px;
  border-top: 1px solid #ebeef5;
}
</style>
