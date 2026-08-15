// Svelte 5 state module for authentication, theme mode, user management, and login history

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

  usersList = $state<RegisteredUser[]>([
    {
      id: '1',
      name: 'Deepak (Admin)',
      email: 'admin@deepphotos.local',
      role: 'Administrator',
      avatar: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80',
      status: 'Active',
      lastLogin: 'Just now'
    },
    {
      id: '2',
      name: 'Sarah Connor',
      email: 'sarah@deepphotos.local',
      role: 'Editor',
      avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150&auto=format&fit=crop&q=80',
      status: 'Active',
      lastLogin: '2 hours ago'
    },
    {
      id: '3',
      name: 'Alex Rivera',
      email: 'alex@deepphotos.local',
      role: 'Viewer',
      avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&auto=format&fit=crop&q=80',
      status: 'Active',
      lastLogin: '1 day ago'
    }
  ]);

  loginHistory = $state<LoginLog[]>([
    { id: '101', user: 'admin@deepphotos.local', timestamp: 'Today at 23:38', ip: '192.168.1.104', device: 'MacBook Pro • macOS (zsh)', status: 'Success' },
    { id: '102', user: 'admin@deepphotos.local', timestamp: 'Today at 21:45', ip: '192.168.1.104', device: 'Safari • macOS', status: 'Success' },
    { id: '103', user: 'sarah@deepphotos.local', timestamp: 'Today at 19:12', ip: '192.168.1.112', device: 'iPhone 15 Pro • iOS', status: 'Success' },
    { id: '104', user: 'unknown@external.io', timestamp: 'Yesterday at 04:10', ip: '45.12.89.201', device: 'Firefox • Linux', status: 'Failed' },
    { id: '105', user: 'alex@deepphotos.local', timestamp: 'Aug 14, 2026 at 14:20', ip: '192.168.1.118', device: 'Chrome • Windows 11', status: 'Success' }
  ]);
  
  isSidebarCollapsed = $state(false);

  login() {
    this.isAuthenticated = true;
  }

  logout() {
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

  addUser(newUser: Omit<RegisteredUser, 'id' | 'lastLogin'>) {
    const user: RegisteredUser = {
      ...newUser,
      id: String(Date.now()),
      lastLogin: 'Never'
    };
    this.usersList.push(user);
  }

  deleteUser(userId: string) {
    this.usersList = this.usersList.filter(u => u.id !== userId);
  }
}

export const appState = new AppState();
