<template>
  <el-container style="height: 100vh">
    <el-aside width="220px" style="background-color: #304156; color: #fff">
      <div class="logo">
        <el-icon :size="24"><OfficeBuilding /></el-icon>
        <span>社团管理系统</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><House /></el-icon>
          <span>首页</span>
        </el-menu-item>

        <el-menu-item
          v-if="userStore.isAdmin()"
          index="/users"
        >
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>

        <el-sub-menu index="clubs">
          <template #title>
            <el-icon><Collection /></el-icon>
            <span>社团管理</span>
          </template>
          <el-menu-item index="/clubs">社团列表</el-menu-item>
          <el-menu-item index="/clubs/my">我的社团</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="activities">
          <template #title>
            <el-icon><Calendar /></el-icon>
            <span>活动管理</span>
          </template>
          <el-menu-item index="/activities">活动列表</el-menu-item>
          <el-menu-item index="/activities/my">我的活动</el-menu-item>
          <el-menu-item index="/activities/organized">我组织的活动</el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/profile">
          <el-icon><UserFilled /></el-icon>
          <span>个人中心</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header style="background-color: #fff; box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08)">
        <div class="header-content">
          <span style="font-size: 16px; font-weight: 500">{{ currentTitle }}</span>
          <div class="user-info">
            <el-dropdown @command="handleCommand">
              <span class="user-dropdown">
                <el-icon :size="20"><UserFilled /></el-icon>
                {{ userStore.user?.name }}
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </el-header>

      <el-main style="background-color: #f0f2f5; padding: 20px">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const currentTitle = computed(() => route.meta.title || '')

const handleCommand = (command) => {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
  border-bottom: 1px solid #1f2d3d;
}

.logo span {
  margin-left: 8px;
}

.header-content {
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #606266;
}

.user-dropdown:hover {
  color: #409EFF;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
