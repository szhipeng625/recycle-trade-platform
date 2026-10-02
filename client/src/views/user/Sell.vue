<template>
  <div>
    <el-card shadow="hover">
      <template #header><span>卖货估价计算器</span></template>
      <el-form :model="form" label-width="120px" style="max-width:600px">
        <el-form-item label="选择货物" required>
          <el-select v-model="form.goods_id" placeholder="选择要卖的货物" style="width:100%" @change="calc">
            <el-option v-for="p in pricing" :key="p.goods_id" :label="p.goods?.name" :value="p.goods_id" />
          </el-select>
        </el-form-item>
        <el-form-item label="今日单价">
          <el-input :value="currentPrice ? '¥' + currentPrice.toFixed(2) + ' /公斤' : '未定价'" disabled />
        </el-form-item>
        <el-form-item label="重量" required>
          <el-input-number v-model="form.weight" :min="0" :precision="3" :step="1" style="width:65%" @input="calc" />
          <el-select v-model="form.unit" style="width:33%;margin-left:2%" @change="calc">
            <el-option label="公斤" value="公斤" />
            <el-option label="市斤" value="市斤" />
            <el-option label="吨" value="吨" />
          </el-select>
        </el-form-item>
        <el-form-item label="折合公斤">
          <el-input :value="weightInKg.toFixed(3) + ' 公斤'" disabled />
        </el-form-item>
        <el-form-item label="预估金额">
          <div class="result">¥{{ totalAmount.toFixed(2) }}</div>
        </el-form-item>
        <el-alert type="info" :closable="false" style="margin-top:10px">
          此为预估金额，实际以现场称重和管理员开单为准。请联系管理员完成交易。
        </el-alert>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import api from '../../api'

const pricing = ref([])
const form = reactive({ goods_id: null, weight: 0, unit: '公斤' })

const currentPrice = computed(() => {
  const p = pricing.value.find(p => p.goods_id === form.goods_id)
  return p ? p.price : 0
})

const weightInKg = computed(() => {
  if (form.unit === '吨') return form.weight * 1000
  if (form.unit === '市斤') return form.weight * 0.5
  return form.weight
})

const totalAmount = computed(() => weightInKg.value * currentPrice.value)

const calc = () => {}

onMounted(async () => {
  const res = await api.get('/pricing/today')
  pricing.value = res.data
})
</script>

<style scoped>
.result { font-size: 32px; font-weight: bold; color: #f56c6c; }
</style>
