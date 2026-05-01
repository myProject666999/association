<template>
  <div>
    <el-card class="common-card">
      <template #header>
        <span>个人信息</span>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="info">
          <el-form
            ref="profileFormRef"
            :model="profileForm"
            :rules="profileRules"
            label-width="100px"
            style="max-width: 500px"
          >
            <el-form-item label="用户名">
              <el-input :value="userInfo?.username" disabled />
            </el-form-item>
            <el-form-item label="姓名" prop="name">
              <el-input v-model="profileForm.name" placeholder="请输入姓名" />
            </el-form-item>
            <el-form-item label="学号" prop="student_id">
              <el-input v-model="profileForm.student_id" placeholder="请输入学号" />
            </el-form-item>
            <el-form-item label="院系" prop="department">
              <el-input v-model="profileForm.department" placeholder="请输入院系" />
            </el-form-item>
            <el-form-item label="专业" prop="major">
              <el-input v-model="profileForm.major" placeholder="请输入专业" />
            </el-form-item>
            <el-form-item label="年级" prop="grade">
              <el-input v-model="profileForm.grade" placeholder="请输入年级" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="profileForm.email" placeholder="请输入邮箱" />
            </el-form-item>
            <el-form-item label="手机" prop="phone">
              <el-input v-model="profileForm.phone" placeholder="请输入手机号" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="submitLoading" @click="handleUpdateProfile">
                保存修改
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="修改密码" name="password">
          <el-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            label-width="100px"
            style="max-width: 500px"
          >
            <el-form-item label="原密码" prop="old_password">
              <el-input
                v-model="passwordForm.old_password"
                type="password"
                placeholder="请输入原密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="新密码" prop="new_password">
              <el-input
                v-model="passwordForm.new_password"
                type="password"
                placeholder="请输入新密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirm_password">
              <el-input
                v-model="passwordForm.confirm_password"
                type="password"
                placeholder="请再次输入新密码"
                show-password
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="passwordLoading" @click="handleUpdatePassword">
                修改密码
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="我的活动" name="activities">
          <el-table :data="myActivities" stripe style="width: 100%" v-loading="activitiesLoading">
            <el-table-column prop="activity.title" label="活动名称" min-width="200">
              <template #default="scope">
                {{ scope.row.activity?.title }}
              </template>
            </el-table-column>
            <el-table-column prop="activity.location" label="地点" width="120">
              <template #default="scope">
                {{ scope.row.activity?.location }}
              </template>
            </el-table-column>
            <el-table-column label="时间" width="180">
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
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="我的社团" name="clubs">
          <el-table :data="myClubs" stripe style="width: 100%" v-loading="clubsLoading">
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
                {{ getPositionText(scope.row.position) }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="120">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)" size="small">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="加入时间" width="180">
              <template #default="scope">
                {{ scope.row.joined_at ? formatDateTime(scope.row.joined_at) : '-' }}
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/utils/request'

const activeTab = ref('info')
const submitLoading = ref(false)
const passwordLoading = ref(false)
const activitiesLoading = ref(false)
const clubsLoading = ref(false)

const userInfo = ref(null)
const myActivities = ref([])
const myClubs = ref([])

const profileFormRef = ref(null)
const profileForm = reactive({
  name: '',
  student_id: '',
  department: '',
  major: '',
  grade: '',
  email: '',
  phone: ''
})

const profileRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}

const passwordFormRef = ref(null)
const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordForm.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  old_password: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const fetchUserInfo = async () => {
  try {
    const res = await api.get('/api/profile')
    if (res.code === 200) {
      userInfo.value = res.data
      profileForm.name = res.data.name || ''
      profileForm.student_id = res.data.student_id || ''
      profileForm.department = res.data.department || ''
      profileForm.major = res.data.major || ''
      profileForm.grade = res.data.grade || ''
      profileForm.email = res.data.email || ''
      profileForm.phone = res.data.phone || ''
    }
  } catch (error) {
    console.error('获取用户信息失败', error)
  }
}

const fetchMyActivities = async () => {
  activitiesLoading.value = true
  try {
    const res = await api.get('/api/profile/activities')
    if (res.code === 200) {
      myActivities.value = res.data || []
    }
  } catch (error) {
    console.error('获取我的活动失败', error)
  } finally {
    activitiesLoading.value = false
  }
}

const fetchMyClubs = async () => {
  clubsLoading.value = true
  try {
    const res = await api.get('/api/profile/clubs')
    if (res.code === 200) {
      myClubs.value = res.data || []
    }
  } catch (error) {
    console.error('获取我的社团失败', error)
  } finally {
    clubsLoading.value = false
  }
}

const handleUpdateProfile = async () => {
  if (!profileFormRef.value) return

  await profileFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        const res = await api.put('/api/profile', profileForm)
        if (res.code === 200) {
          ElMessage.success('个人信息更新成功')
          fetchUserInfo()
        }
      } catch (error) {
        console.error('更新失败', error)
      } finally {
        submitLoading.value = false
      }
    }
  })
}

const handleUpdatePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      passwordLoading.value = true
      try {
        const data = {
          old_password: passwordForm.old_password,
          new_password: passwordForm.new_password
        }
        const res = await api.put('/api/profile/password', data)
        if (res.code === 200) {
          ElMessage.success('密码修改成功')
          passwordForm.old_password = ''
          passwordForm.new_password = ''
          passwordForm.confirm_password = ''
        }
      } catch (error) {
        console.error('修改密码失败', error)
      } finally {
        passwordLoading.value = false
      }
    }
  })
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

const getPositionText = (position) => {
  switch (position) {
    case 'president': return '社长'
    case 'vice_president': return '副社长'
    default: return '成员'
  }
}

onMounted(() => {
  fetchUserInfo()
  fetchMyActivities()
  fetchMyClubs()
})
</script>

<style scoped>
</style>
