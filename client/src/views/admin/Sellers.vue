<template>
  <div>
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>卖货人管理</span>
          <el-button type="primary" @click="openDialog">新增卖货人</el-button>
        </div>
      </template>
      <el-input v-model="keyword" placeholder="搜索姓名或手机号" style="width:300px;margin-bottom:16px" clearable @change="loadData" />
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" width="150" />
        <el-table-column prop="phone" label="手机号" width="150" />
        <el-table-column prop="id_card" label="身份证号" width="200" />
        <el-table-column prop="address" label="地址" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑卖货人' : '新增卖货人'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="姓名" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="身份证号">
          <el-input v-model="form.id_card" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="form.address" type="textarea" :rows="2" />
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
const keyword = ref('')
const dialogVisible = ref(false)
const editingId = ref(null)
const form = reactive({ name: '', phone: '', id_card: '', address: '' })

const loadData = async () => {
  const res = await api.get('/sellers', { params: { keyword: keyword.value } })
  list.value = res.data
}

const openDialog = (row) => {
  if (row) {
    editingId.value = row.id
    Object.assign(form, row)
  } else {
    editingId.value = null
    form.name = ''; form.phone = ''; form.id_card = ''; form.address = ''
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!form.name) { ElMessage.warning('请输入姓名'); return }
  if (editingId.value) {
    await api.put(`/sellers/${editingId.value}`, form)
  } else {
    await api.post('/sellers', form)
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  loadData()
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定删除 "${row.name}"？`, '提示', { type: 'warning' })
    .then(async () => {
      await api.delete(`/sellers/${row.id}`)
      ElMessage.success('删除成功')
      loadData()
    }).catch(() => {})
}

onMounted(loadData)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
