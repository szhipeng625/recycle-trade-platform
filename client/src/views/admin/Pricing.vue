<template>
  <div>
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>每日定价管理</span>
          <el-button type="primary" @click="openDialog">新增定价</el-button>
        </div>
      </template>
      <el-form :inline="true" style="margin-bottom:16px">
        <el-form-item label="日期">
          <el-date-picker v-model="filterDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" @change="loadData" />
        </el-form-item>
        <el-form-item label="货物">
          <el-select v-model="filterGoods" placeholder="全部" clearable style="width:150px" @change="loadData">
            <el-option v-for="g in goodsList" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="goods.name" label="货物" width="150" />
        <el-table-column prop="goods.category" label="类别" width="150" />
        <el-table-column label="单价(元/公斤)" width="160">
          <template #default="{ row }"><span style="color:#f56c6c;font-weight:bold;font-size:16px">¥{{ row.price.toFixed(2) }}</span></template>
        </el-table-column>
        <el-table-column prop="date" label="定价日期" width="150" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新增/修改定价" width="450px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="货物" required>
          <el-select v-model="form.goods_id" placeholder="请选择货物" style="width:100%">
            <el-option v-for="g in goodsList" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="单价(元/公斤)" required>
          <el-input-number v-model="form.price" :min="0" :precision="2" :step="0.1" style="width:100%" />
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker v-model="form.date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
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
const goodsList = ref([])
const filterDate = ref('')
const filterGoods = ref('')
const dialogVisible = ref(false)
const form = reactive({ goods_id: null, price: 0, date: '' })

const loadData = async () => {
  const params = {}
  if (filterDate.value) params.date = filterDate.value
  if (filterGoods.value) params.goods_id = filterGoods.value
  const res = await api.get('/pricing', { params })
  list.value = res.data
}

const loadGoods = async () => {
  const res = await api.get('/goods')
  goodsList.value = res.data
}

const openDialog = () => {
  form.goods_id = null
  form.price = 0
  form.date = new Date().toISOString().slice(0, 10)
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!form.goods_id) { ElMessage.warning('请选择货物'); return }
  if (form.price <= 0) { ElMessage.warning('请输入单价'); return }
  await api.post('/pricing', form)
  ElMessage.success('保存成功')
  dialogVisible.value = false
  loadData()
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定删除 ${row.goods?.name} 的定价记录？`, '提示', { type: 'warning' })
    .then(async () => {
      await api.delete(`/pricing/${row.id}`)
      ElMessage.success('删除成功')
      loadData()
    }).catch(() => {})
}

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : ''

onMounted(() => { loadGoods(); loadData() })
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
