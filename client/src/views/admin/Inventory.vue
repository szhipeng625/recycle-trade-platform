<template>
  <div>
    <el-row :gutter="16" style="margin-bottom:16px">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>库存总览</span>
              <el-button type="primary" size="small" @click="openInDialog">入库</el-button>
              <el-button type="warning" size="small" @click="openOutDialog">出库</el-button>
            </div>
          </template>
          <el-table :data="inventory" stripe>
            <el-table-column prop="goods.name" label="货物" width="150" />
            <el-table-column prop="goods.category" label="类别" width="150" />
            <el-table-column label="库存(公斤)" width="150">
              <template #default="{ row }">
                <span :style="{color: row.quantity < 100 ? '#f56c6c' : '#67c23a', fontWeight:'bold'}">{{ row.quantity.toFixed(2) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="更新时间">
              <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header><span>出入库流水</span></template>
          <el-table :data="records" stripe size="small" max-height="400">
            <el-table-column label="类型" width="80">
              <template #default="{ row }">
                <el-tag :type="row.type === 'in' ? 'success' : 'danger'" size="small">{{ row.type === 'in' ? '入库' : '出库' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="货物" width="100">
              <template #default="{ row }">{{ row.goods?.name }}</template>
            </el-table-column>
            <el-table-column label="数量(公斤)" width="110">
              <template #default="{ row }">{{ row.quantity.toFixed(2) }}</template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" show-overflow-tooltip />
            <el-table-column label="时间" width="150">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="inDialog" title="入库" width="450px">
      <el-form :model="inForm" label-width="100px">
        <el-form-item label="货物" required>
          <el-select v-model="inForm.goods_id" style="width:100%">
            <el-option v-for="g in goodsList" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="数量" required>
          <el-input-number v-model="inForm.quantity" :min="0" :precision="3" style="width:65%" />
          <el-select v-model="inForm.unit" style="width:33%;margin-left:2%">
            <el-option label="公斤" value="公斤" />
            <el-option label="市斤" value="市斤" />
            <el-option label="吨" value="吨" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="inForm.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="inDialog = false">取消</el-button>
        <el-button type="primary" @click="handleIn">确认入库</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="outDialog" title="出库" width="450px">
      <el-form :model="outForm" label-width="100px">
        <el-form-item label="货物" required>
          <el-select v-model="outForm.goods_id" style="width:100%">
            <el-option v-for="g in goodsList" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="数量" required>
          <el-input-number v-model="outForm.quantity" :min="0" :precision="3" style="width:65%" />
          <el-select v-model="outForm.unit" style="width:33%;margin-left:2%">
            <el-option label="公斤" value="公斤" />
            <el-option label="市斤" value="市斤" />
            <el-option label="吨" value="吨" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="outForm.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="outDialog = false">取消</el-button>
        <el-button type="warning" @click="handleOut">确认出库</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../../api'

const inventory = ref([])
const records = ref([])
const goodsList = ref([])
const inDialog = ref(false)
const outDialog = ref(false)
const inForm = reactive({ goods_id: null, quantity: 0, unit: '公斤', remark: '' })
const outForm = reactive({ goods_id: null, quantity: 0, unit: '公斤', remark: '' })

const loadData = async () => {
  const [inv, rec, goods] = await Promise.all([
    api.get('/inventory'),
    api.get('/inventory/records', { params: { size: 50 } }),
    api.get('/goods')
  ])
  inventory.value = inv.data
  records.value = rec.data.list
  goodsList.value = goods.data
}

const openInDialog = () => {
  inForm.goods_id = null; inForm.quantity = 0; inForm.unit = '公斤'; inForm.remark = ''
  inDialog.value = true
}

const openOutDialog = () => {
  outForm.goods_id = null; outForm.quantity = 0; outForm.unit = '公斤'; outForm.remark = ''
  outDialog.value = true
}

const handleIn = async () => {
  if (!inForm.goods_id || inForm.quantity <= 0) { ElMessage.warning('请填写完整'); return }
  await api.post('/inventory/in', inForm)
  ElMessage.success('入库成功')
  inDialog.value = false
  loadData()
}

const handleOut = async () => {
  if (!outForm.goods_id || outForm.quantity <= 0) { ElMessage.warning('请填写完整'); return }
  try {
    await api.post('/inventory/out', outForm)
    ElMessage.success('出库成功')
    outDialog.value = false
    loadData()
  } catch (e) {}
}

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : ''

onMounted(loadData)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; gap: 8px; }
</style>
