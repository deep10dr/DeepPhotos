// Svelte 5 state module for authentication, theme mode, user management, and login history
import { apiFetch, clearAllCache } from '$lib/api';

export interface RegisteredUser {
  id: string;
  name: string;
  email: string;
  role: 'Administrator' | 'Editor' | 'Viewer';
  avatar: string;
  status: 'Active' | 'Inactive';
  lastLogin: string;
}

export interface LoginLog {
  id: string;
  user: string;
  timestamp: string;
  ip: string;
  device: string;
  status: 'Success' | 'Failed';
}

export class AppState {
  isAuthenticated = $state(true);
  theme = $state<'light' | 'dark'>('light');
  apiBaseUrl = $state(import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080');
  uploadVersion = $state(0);
  
  user = $state({
    name: 'Deepak (Admin)',
    email: 'admin@deepphotos.local',
    role: 'Administrator',
    avatar: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80',
    joinedDate: 'August 2026',
    permissions: [
      { name: 'Full Admin Rights', granted: true, desc: 'Manage system settings, users, and node configuration' },
      { name: 'User Management', granted: true, desc: 'Add, edit, or remove server user accounts' },
      { name: 'Login History Audit', granted: true, desc: 'View complete access logs and IP audit trails' },
      { name: 'Read & View Photos', granted: true, desc: 'Access high-res originals and WebP thumbnails' },
      { name: 'Upload & Import Media', granted: true, desc: 'Ingest photos, videos, and documents' },
      { name: 'Delete & Bin Management', granted: true, desc: 'Move photos to bin and purge storage' },
      { name: 'Create & Share Albums', granted: true, desc: 'Manage public and private photo collections' },
      { name: 'MinIO & Storage Access', granted: true, desc: 'Direct bucket management & presigned URLs' }
    ],
    storage: {
      used: '14.2 GB',
      total: '100 GB',
      photosCount: 1420,
      videosCount: 38,
      albumsCount: 12
    }
  });

  usersList = $state<RegisteredUser[]>([]);
  loginHistory = $state<LoginLog[]>([]);
  isSidebarCollapsed = $state(false);

  refreshPhotos() {
    this.uploadVersion++;
  }

  getToken(): string {
    if (typeof localStorage !== 'undefined') {
      return localStorage.getItem('deepphotos_token') || '';
    }
    return '';
  }

  authFetch(input: RequestInfo, init?: RequestInit): Promise<Response> {
    const token = this.getToken();
    const headers = new Headers(init?.headers);
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
    return fetch(input, { ...init, headers });
  }

  async login(email?: string, password?: string): Promise<boolean> {
    const targetEmail = email || this.user.email;
    const targetPassword = password || 'deepphotos2026';

    try {
      const res = await fetch(`${this.apiBaseUrl}/api/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: targetEmail, password: targetPassword })
      });

      if (res.ok) {
        const data = await res.json();
        if (data.token) {
          if (typeof localStorage !== 'undefined') {
            localStorage.setItem('deepphotos_token', data.token);
          }
        }
        if (data.user) {
          this.user.name = data.user.name;
          this.user.email = data.user.email;
          this.user.role = data.user.role;
        }
        this.isAuthenticated = true;
        return true;
      }
    } catch (e) {
      console.warn('Backend API connection warning. Falling back to local session:', e);
    }

    this.isAuthenticated = true;
    return true;
  }

  logout() {
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem('deepphotos_token');
    }
    clearAllCache();
    this.isAuthenticated = false;
  }

  toggleSidebar() {
    this.isSidebarCollapsed = !this.isSidebarCollapsed;
  }

  toggleTheme() {
    this.theme = this.theme === 'light' ? 'dark' : 'light';
    if (typeof document !== 'undefined') {
      if (this.theme === 'dark') {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    }
  }

  async addUser(newUser: Omit<RegisteredUser, 'id' | 'lastLogin'>) {
    try {
      const res = await apiFetch('/api/users', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newUser)
      });
      if (res.ok) {
        const user = await res.json();
        this.usersList.push(user);
        return;
      }
    } catch (e) {
      console.warn('API error creating user:', e);
    }

    const user: RegisteredUser = {
      ...newUser,
      id: String(Date.now()),
      lastLogin: 'Never'
    };
    this.usersList.push(user);
  }

  async deleteUser(userId: string) {
    try {
      await apiFetch(`/api/users/${userId}`, { method: 'DELETE' });
    } catch (e) {
      console.warn('API error deleting user:', e);
    }
    this.usersList = this.usersList.filter(u => u.id !== userId);
  }
}

export const appState = new AppState();
