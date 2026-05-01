<template>
  <div>
    <el-card class="common-card">
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新增用户
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="姓名/用户名/学号" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="searchForm.role" placeholder="全部角色" clearable style="width: 150px">
            <el-option label="校级管理员" value="university_admin" />
            <el-option label="院级管理员" value="dept_admin" />
            <el-option label="学生" value="student" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="userList" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="student_id" label="学号" width="120" />
        <el-table-column prop="department" label="院系" width="150" />
        <el-table-column prop="major" label="专业" width="120" />
        <el-table-column prop="role" label="角色" width="120">
          <template #default="scope">
            <el-tag :type="getRoleTagType(scope.row.role)">
              {{ getRoleText(scope.row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-switch
              v-model="scope.row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleToggleStatus(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="scope">
            <el-button type="primary" link size="small" @click="handleEdit(scope.row)">
              编辑
            </el-button>
            <el-button type="warning" link size="small" @click="handleChangeRole(scope.row)">
              分配角色
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="userFormRef" :model="userForm" :rules="userRules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" placeholder="请输入用户名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input v-model="userForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="userForm.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="学号" prop="student_id">
          <el-input v-model="userForm.student_id" placeholder="请输入学号" />
        </el-form-item>
        <el-form-item label="院系" prop="department">
          <el-input v-model="userForm.department" placeholder="请输入院系" />
        </el-form-item>
        <el-form-item label="专业" prop="major">
          <el-input v-model="userForm.major" placeholder="请输入专业" />
        </el-form-item>
        <el-form-item label="年级" prop="grade">
          <el-input v-model="userForm.grade" placeholder="请输入年级" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机" prop="phone">
          <el-input v-model="userForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="校级管理员" value="university_admin" />
            <el-option label="院级管理员" value="dept_admin" />
            <el-option label="学生" value="student" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="roleDialogVisible" title="分配角色" width="400px">
      <el-form label-width="80px">
        <el-form-item label="当前用户">
          <el-input :value="currentUser?.name" disabled />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="roleForm.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="校级管理员" value="university_admin" />
            <el-option label="院级管理员" value="dept_admin" />
            <el-option label="学生" value="student" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitRole">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const userList = ref([])
const dialogVisible = ref(false)
const roleDialogVisible = ref(false)
const isEdit = ref(false)
const currentUser = ref(null)

const searchForm = reactive({
  keyword: '',
  role: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const userFormRef = ref(null)
const userForm = reactive({
  id: null,
  username: '',
  password: '',
  name: '',
  student_id: '',
  department: '',
  major: '',
  grade: '',
  email: '',
  phone: '',
  role: 'student'
})

const roleForm = reactive({
  role: ''
})

const userRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const dialogTitle = ref('新增用户')

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.role) params.role = searchForm.role

    const res = await api.get('/api/users', { params })
    if (res.code === 200) {
      userList.value = res.data.list
      pagination.total = res.data.total
    }
  } catch (error) {
    console.error('获取用户列表失败', error)
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
  searchForm.role = ''
  pagination.page = 1
  fetchData()
}

const handleCreate = () => {
  isEdit.value = false
  dialogTitle.value = '新增用户'
  userForm.id = null
  userForm.username = ''
  userForm.password = ''
  userForm.name = ''
  userForm.student_id = ''
  userForm.department = ''
  userForm.major = ''
  userForm.grade = ''
  userForm.email = ''
  userForm.phone = ''
  userForm.role = 'student'
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑用户'
  userForm.id = row.id
  userForm.username = row.username
  userForm.name = row.name
  userForm.student_id = row.student_id
  userForm.department = row.department
  userForm.major = row.major
  userForm.grade = row.grade
  userForm.email = row.email
  userForm.phone = row.phone
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!userFormRef.value) return

  await userFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (isEdit.value) {
          const data = { ...userForm }
          delete data.password
          const res = await api.put(`/api/users/${userForm.id}`, data)
          if (res.code === 200) {
            ElMessage.success('更新成功')
            dialogVisible.value = false
            fetchData()
          }
        } else {
          const res = await api.post('/api/users', userForm)
          if (res.code === 200) {
            ElMessage.success('创建成功')
            dialogVisible.value = false
            fetchData()
          }
        }
      } catch (error) {
        console.error('提交失败', error)
      } finally {
        submitLoading.value = false
      }
    }
  })
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该用户吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await api.delete(`/api/users/${row.id}`)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        fetchData()
      }
    } catch (error) {
      console.error('删除失败', error)
    }
  }).catch(() => {})
}

const handleToggleStatus = async (row) => {
  try {
    const res = await api.put(`/api/users/${row.id}/status`)
    if (res.code === 200) {
      ElMessage.success('状态更新成功')
    }
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1
    console.error('状态更新失败', error)
  }
}

const handleChangeRole = (row) => {
  currentUser.value = row
  roleForm.role = row.role
  roleDialogVisible.value = true
}

const handleSubmitRole = async () => {
  try {
    const res = await api.put(`/api/users/${currentUser.value.id}/role`, { role: roleForm.role })
    if (res.code === 200) {
      ElMessage.success('角色分配成功')
      roleDialogVisible.value = false
      fetchData()
    }
  } catch (error) {
    console.error('角色分配失败', error)
  }
}

const getRoleTagType = (role) => {
  switch (role) {
    case 'university_admin': return 'danger'
    case 'dept_admin': return 'warning'
    case 'student': return 'success'
    default: return 'info'
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
