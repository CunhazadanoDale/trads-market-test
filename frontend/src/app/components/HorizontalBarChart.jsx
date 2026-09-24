import React from 'react';
import { Bar, BarChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';
import { ChartTooltip } from './ChartTooltip';
import './HorizontalBarChart.css';

const LABEL_AXIS_ID = 'labels';

export function HorizontalBarChart({
  data,
  labelKey,
  barKey,
  formatValue,
  tooltip,
  valueAxis,
  rowHeight = 36,
  barSize = 18,
  labelWidth = 110,
  valueWidth = 96,
}) {
  if (!data || data.length === 0) return null;

  const handleTooltip = ({ active, payload }) => {
    const item = payload && payload.length > 0 ? payload[0].payload : null;
    if (!active || !item) return null;
    const content = tooltip(item);
    return <ChartTooltip title={content.title} rows={content.rows} />;
  };

  const handleLabel = ({ x, y, width, height, value }) => {
    if (value == null) return null;
    return (
      <text className="hbar-value" x={x + width + 8} y={y + height / 2} dy="0.35em">
        {formatValue(value)}
      </text>
    );
  };

  return (
    <div className="hbar-chart">
      <ResponsiveContainer width="100%" height={data.length * rowHeight}>
        <BarChart
          data={data}
          layout="vertical"
          margin={{ top: 0, right: valueAxis ? 8 : valueWidth, bottom: 0, left: 0 }}
          barCategoryGap="24%"
        >
          <XAxis type="number" hide dataKey={barKey} />
          <YAxis
            type="category"
            yAxisId={LABEL_AXIS_ID}
            dataKey={labelKey}
            width={labelWidth}
            interval={0}
            tickLine={false}
            axisLine={false}
          />
          {valueAxis ? (
            <YAxis
              type="category"
              yAxisId="values"
              orientation="right"
              dataKey={valueAxis.dataKey}
              width={valueAxis.width}
              interval={0}
              allowDuplicatedCategory={false}
              tickLine={false}
              axisLine={false}
              tickFormatter={valueAxis.format}
            />
          ) : null}
          <Tooltip axisId={LABEL_AXIS_ID} cursor={{ className: 'hbar-cursor' }} content={handleTooltip} />
          <Bar
            dataKey={barKey}
            yAxisId={LABEL_AXIS_ID}
            barSize={barSize}
            radius={[0, 4, 4, 0]}
            label={formatValue ? handleLabel : undefined}
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
