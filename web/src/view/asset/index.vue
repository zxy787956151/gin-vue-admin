<template>
  <div class="asset-view-container">
    <el-card class="box-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">个人资产视图</span>
          <el-button type="primary" :icon="Refresh" circle @click="fetchData" :loading="loading" />
        </div>
      </template>
      
      <div v-loading="loading" class="content-wrapper">
        <div class="summary-section">
          <el-statistic title="总资产" :value="totalAsset" :precision="2">
            <template #prefix>
              <el-icon style="vertical-align: -0.125em"><Money /></el-icon>
            </template>
            <template #suffix>元</template>
          </el-statistic>
        </div>
        
        <div class="chart-section" v-if="assetData.items.length > 0">
          <Chart 
            :height="'500px'" 
            :option="chartOption" 
          />
        </div>
        <div v-else class="chart-section empty-chart">
          <el-empty description="暂无数据" />
        </div>
        
        <div class="table-section">
          <el-table :data="tableData" stripe style="width: 100%">
            <el-table-column prop="name" label="资产类型" width="180" />
            <el-table-column prop="value" label="金额" width="180">
              <template #default="{ row }">
                <span>{{ formatCurrency(row.value) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="percentage" label="占比">
              <template #default="{ row }">
                <el-progress :percentage="row.percentage" :color="getProgressColor(row.percentage)" />
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Money } from '@element-plus/icons-vue'
import Chart from '@/components/charts/index.vue'
import useChartOption from '@/hooks/charts'
import { useAppStore } from '@/pinia'
import { getAssetDistribution } from '@/api/asset'

defineOptions({
  name: 'AssetView'
})

const loading = ref(false)
const assetData = ref({
  total: 0,
  items: [],
  timestamp: 0
})

// 计算总资产
const totalAsset = computed(() => {
  return assetData.value.total
})

// 计算表格数据（包含占比）
const tableData = computed(() => {
  return assetData.value.items.map(item => ({
    ...item,
    percentage: assetData.value.total > 0 
      ? parseFloat(((item.value / assetData.value.total) * 100).toFixed(2))
      : 0
  }))
})

// 饼图配置
const colors = [
  '#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc'
]

// 使用useChartOption，但通过访问assetData.value确保响应式追踪
const { chartOption } = useChartOption((isDark) => {
  // 在函数内部访问assetData.value，useChartOption的computed会追踪这个依赖
  const pieData = assetData.value.items.map((item, index) => ({
    value: item.value,
    name: item.name,
    itemStyle: {
      color: colors[index % colors.length]
    }
  }))
  
  return {
    tooltip: {
      trigger: 'item',
      formatter: (params) => {
        const percentage = assetData.value.total > 0 
          ? ((params.value / assetData.value.total) * 100).toFixed(2)
          : 0
        return `${params.name}<br/>金额: ${formatCurrency(params.value)}<br/>占比: ${percentage}%`
      }
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      top: 'center'
    },
    series: [
      {
        name: '资产分布',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: isDark ? '#1e1e1e' : '#fff',
          borderWidth: 2
        },
        label: {
          show: true,
          formatter: (params) => {
            const percentage = assetData.value.total > 0
              ? ((params.value / assetData.value.total) * 100).toFixed(1)
              : 0
            return `${params.name}\n${percentage}%`
          }
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          },
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        data: pieData
      }
    ]
  }
})

// 格式化货币
const formatCurrency = (value) => {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY',
    minimumFractionDigits: 2
  }).format(value)
}

// 获取进度条颜色
const getProgressColor = (percentage) => {
  if (percentage >= 30) return '#67c23a'
  if (percentage >= 15) return '#e6a23c'
  return '#f56c6c'
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getAssetDistribution()
    if (res.code === 0) {
      assetData.value = res.data
      ElMessage.success('数据加载成功')
    } else {
      ElMessage.error(res.msg || '获取数据失败')
    }
  } catch (error) {
    console.error('获取资产数据失败:', error)
    ElMessage.error('获取数据失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 页面挂载时获取数据
onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.asset-view-container {
  padding: 20px;
  
  .box-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .title {
        font-size: 18px;
        font-weight: bold;
      }
    }
    
    .content-wrapper {
      .summary-section {
        margin-bottom: 30px;
        text-align: center;
        padding: 20px;
        background: var(--el-bg-color-page);
        border-radius: 8px;
      }
      
      .chart-section {
        margin-bottom: 30px;
        padding: 20px;
        background: var(--el-bg-color-page);
        border-radius: 8px;
      }
      
      .table-section {
        margin-top: 20px;
      }
    }
  }
}
</style>
