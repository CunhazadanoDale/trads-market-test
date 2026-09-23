import React, { useState } from 'react';
import { NavLink } from 'react-router-dom';
import {
  LayoutDashboard,
  Map as MapIcon,
  MapPin,
  Users,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';
import './Sidebar.css';

export function Sidebar() {
  const [collapsed, setCollapsed] = useState(false);

  // Menu espelhando os recursos disponíveis no backend-golang
  const navGroups = [
    {
      title: 'Geral',
      items: [
        { label: 'Dashboard', icon: <LayoutDashboard size={18} />, path: '/' },
      ],
    },
    {
      title: 'Inteligência de Mercado',
      items: [
        { label: 'Público', icon: <Users size={18} />, path: '/publico' },
      ],
    },
    {
      title: 'Localização (IBGE)',
      items: [
        { label: 'Estados', icon: <MapIcon size={18} />, path: '/estados' },
        { label: 'Cidades', icon: <MapPin size={18} />, path: '/cidades' },
      ],
    },
  ];

  return (
    <aside className={`sidebar ${collapsed ? 'collapsed' : ''}`}>
      <div className="sidebar-header">
        {collapsed ? 'T' : 'TRADS'}
      </div>

      <div className="sidebar-content">
        {navGroups.map((group) => (
          <div key={group.title} className="sidebar-group">
            <div className="sidebar-group-title">
              {collapsed ? '...' : group.title}
            </div>
            <ul className="sidebar-nav">
              {group.items.map((item) => (
                <li key={item.path}>
                  <NavLink
                    to={item.path}
                    end={item.path === '/'}
                    className={({ isActive }) => (isActive ? 'sidebar-link active' : 'sidebar-link')}
                  >
                    <span className="sidebar-icon">{item.icon}</span>
                    <span className="sidebar-label">{item.label}</span>
                  </NavLink>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>

      <div className="sidebar-footer">
        <button
          type="button"
          className="toggle-btn"
          onClick={() => setCollapsed(!collapsed)}
          title={collapsed ? 'Expandir menu' : 'Recolher menu'}
        >
          {collapsed ? <ChevronRight size={16} /> : <ChevronLeft size={16} />}
        </button>
      </div>
    </aside>
  );
}
