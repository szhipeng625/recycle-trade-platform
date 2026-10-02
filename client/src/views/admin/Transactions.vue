<template>
  <div>
    <el-card shadow="hover" style="margin-bottom:16px">
      <template #header>
        <div class="card-header">
          <span>开单结算</span>
          <el-tag type="success">今日定价自动取价</el-tag>
        </div>
      </template>
      <el-form :model="form" :inline="false" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="卖货人" required>
              <el-select v-model="form.seller_id" placeholder="选择卖货人" filterable allow-create default-first-option style="width:100%" @change="onSellerChange">
                <el-option v-for="s in sellers" :key="s.id" :label="s.name + (s.phone ? ' (' + s.phone + ')' : '')" :value="s.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="货物" required>
              <el-select v-model="form.goods_id" placeholder="选择货物" style="width:100%" @change="calcTotal">
                <el-option v-for="g in goodsList" :key="g.id" :label="g.name" :value="g.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="重量" required>
              <el-input-number v-model="form.weight" :min="0" :precision="3" :step="1" style="width:70%" />
              <el-select v-model="form.unit" style="width:28%;margin-left:2%" @change="calcTotal">
                <el-option label="公斤" value="公斤" />
                <el-option label="市斤" value="市斤" />
                <el-option label="吨" value="吨" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="今日单价">
              <el-input :value="currentPrice ? '¥' + currentPrice.toFixed(2) + ' /公斤' : '未定价'" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="折合公斤">
              <el-input :value="weightInKg.toFixed(3) + ' 公斤'" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="应付金额">
              <el-input :value="'¥' + totalAmount.toFixed(2)" disabled>
                <template #prepend>合计</template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <el-button type="primary" size="large" :disabled="!canSubmit" @click="handleSubmit">确认开单</el-button>
          <el-button size="large" @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>交易记录</span>
          <el-form :inline="true">
            <el-form-item><el-date-picker v-model="filter.start" type="date" value-format="YYYY-MM-DD" placeholder="开始日期" @change="loadList" /></el-form-item>
            <el-form-item><el-date-picker v-model="filter.end" type="date" value-format="YYYY-MM-DD" placeholder="结束日期" @change="loadList" /></el-form-item>
          </el-form>
        </div>
      </template>
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="单号" width="80" />
        <el-table-column label="卖货人" width="120">
          <template #default="{ row }">{{ row.seller?.name }}</template>
        </el-table-column>
        <el-table-column label="货物" width="100">
          <template #default="{ row }">{{ row.goods?.name }}</template>
        </el-table-column>
        <el-table-column label="重量" width="140">
          <template #default="{ row }">{{ row.weight }} {{ row.unit }} ({{ row.weight_in_kg }}公斤)</template>
        </el-table-column>
        <el-table-column label="单价" width="120">
          <template #default="{ row }">¥{{ row.unit_price.toFixed(2) }}/公斤</template>
        </el-table-column>
        <el-table-column label="金额" width="120">
          <template #default="{ row }"><span style="color:#f56c6c;font-weight:bold">¥{{ row.total_amount.toFixed(2) }}</span></template>
        </el-table-column>
        <el-table-column label="操作员" width="100">
          <template #default="{ row }">{{ row.operator?.name || row.operator?.username }}</template>
        </el-table-column>
        <el-table-column label="时间" width="170">
          <template #default="{ row }">{{ formatTime(row.transaction_time) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination style="margin-top:16px;text-align:right" background layout="total, prev, pager, next" :total="total" :page-size="20" :current-page="page" @current-change="p => { page = p; loadList() }" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../../api'

const sellers = ref([])
const goodsList = ref([])
const todayPricing = ref([])
const list = ref([])
const total = ref(0)
const page = ref(1)
const filter = reactive({ start: '', end: '' })

const form = reactive({ seller_id: null, goods_id: null, weight: 0, unit: '公斤' })

const currentPrice = computed(() => {
  const p = todayPricing.value.find(p => p.goods_id === form.goods_id)
  return p ? p.price : 0
})

const weightInKg = computed(() => {
  if (form.unit === '吨') return form.weight * 1000
  if (form.unit === '市斤') return form.weight * 0.5
  return form.weight
})

const totalAmount = computed(() => weightInKg.value * currentPrice.value)

const canSubmit = computed(() => form.seller_id && form.goods_id && form.weight > 0 && currentPrice.value > 0)

const onSellerChange = () => {}

const calcTotal = () => {}

const loadSellers = async () => {
  const res = await api.get('/sellers')
  sellers.value = res.data
}

const loadGoods = async () => {
  const res = await api.get('/goods')
  goodsList.value = res.data
}

const loadPricing = async () => {
  const res = await api.get('/pricing/today')
  todayPricing.value = res.data
}

const loadList = async () => {
  const params = { page: page.value, size: 20 }
  if (filter.start) params.start_date = filter.start
  if (filter.end) params.end_date = filter.end
  const res = await api.get('/transactions', { params })
  list.value = res.data.list
  total.value = res.data.total
}

const handleSubmit = async () => {
  try {
    await api.post('/transactions', form)
    ElMessage.success(`开单成功，应付 ¥${totalAmount.value.toFixed(2)}`)
    resetForm()
    loadList()
  } catch (e) {}
}

const resetForm = () => {
  form.seller_id = null
  form.goods_id = null
  form.weight = 0
  form.unit = '公斤'
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定删除交易 #${row.id}？库存将同步回滚`, '提示', { type: 'warning' })
    .then(async () => {
      await api.delete(`/transactions/${row.id}`)
      ElMessage.success('删除成功')
      loadList()
    }).catch(() => {})
}

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : ''

onMounted(() => {
  loadSellers()
  loadGoods()
  loadPricing()
  loadList()
})
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
