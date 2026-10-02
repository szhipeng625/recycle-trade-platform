<template>
  <div>
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>货物管理</span>
          <el-button type="primary" @click="openDialog">新增货物</el-button>
        </div>
      </template>
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="货物名称" width="200" />
        <el-table-column prop="category" label="类别" width="200" />
        <el-table-column prop="unit" label="默认单位" width="150" />
        <el-table-column prop="created_at" label="创建时间">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑货物' : '新增货物'" width="450px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="如：废铁、废铜" />
        </el-form-item>
        <el-form-item label="类别">
          <el-input v-model="form.category" placeholder="如：黑色金属、有色金属" />
        </el-form-item>
        <el-form-item label="默认单位">
          <el-select v-model="form.unit" style="width:100%">
            <el-option label="公斤" value="公斤" />
            <el-option label="市斤" value="市斤" />
            <el-option label="吨" value="吨" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../../api'

const list = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const form = reactive({ name: '', category: '', unit: '公斤' })

const loadData = async () => {
  const res = await api.get('/goods')
  list.value = res.data
}

const openDialog = (row) => {
  if (row) {
    editingId.value = row.id
    form.name = row.name
    form.category = row.category
    form.unit = row.unit
  } else {
    editingId.value = null
    form.name = ''
    form.category = ''
    form.unit = '公斤'
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!form.name) { ElMessage.warning('请输入名称'); return }
  if (editingId.value) {
    await api.put(`/goods/${editingId.value}`, form)
  } else {
    await api.post('/goods', form)
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  loadData()
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定删除 "${row.name}"？`, '提示', { type: 'warning' })
    .then(async () => {
      await api.delete(`/goods/${row.id}`)
      ElMessage.success('删除成功')
      loadData()
    }).catch(() => {})
}

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : ''

onMounted(loadData)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
