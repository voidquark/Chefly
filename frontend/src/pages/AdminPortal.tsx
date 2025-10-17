import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { usePageTitle } from '../hooks/usePageTitle';
import { apiClient } from '../api/client';
import { ConfirmDialog } from '../components/ConfirmDialog';
import { toast } from 'react-hot-toast';
import type { AdminStats, UserWithStats } from '../types';

export const AdminPortal: React.FC = () => {
  usePageTitle('Admin Portal');
  const { user } = useAuth();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [users, setUsers] = useState<UserWithStats[]>([]);
  const [selectedUser, setSelectedUser] = useState<UserWithStats | null>(null);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [editingLimitUserId, setEditingLimitUserId] = useState<string | null>(null);
  const [limitValue, setLimitValue] = useState<string>('');
  const [showCustomInput, setShowCustomInput] = useState(false);
  const [customLimitValue, setCustomLimitValue] = useState<string>('');

  useEffect(() => {
    // Redirect if not admin
    if (!user?.is_admin) {
      navigate('/dashboard');
      toast.error('Admin access required');
      return;
    }

    fetchData();
  }, [user, navigate]);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [statsData, usersData] = await Promise.all([
        apiClient.getAdminStats(),
        apiClient.getAllUsers(),
      ]);
      setStats(statsData);
      setUsers(usersData.users);
    } catch (error) {
      toast.error('Failed to load admin data');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteClick = (userToDelete: UserWithStats) => {
    if (userToDelete.id === user?.id) {
      toast.error('Cannot delete your own admin account');
      return;
    }
    setSelectedUser(userToDelete);
    setShowDeleteConfirm(true);
  };

  const handleDeleteConfirm = async () => {
    if (!selectedUser) return;

    try {
      setDeleting(true);
      await apiClient.deleteUser(selectedUser.id);
      toast.success(`User ${selectedUser.username} deleted successfully`);
      setShowDeleteConfirm(false);
      setSelectedUser(null);
      // Refresh data
      fetchData();
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to delete user');
    } finally {
      setDeleting(false);
    }
  };

  const handleEditLimit = (userItem: UserWithStats) => {
    setEditingLimitUserId(userItem.id);
    setShowCustomInput(false);
    setCustomLimitValue('');

    // Set initial value: null/undefined -> 'global', -1 -> 'unlimited', 0 -> '0', >0 -> number
    if (userItem.recipe_limit === null || userItem.recipe_limit === undefined) {
      setLimitValue('global');
    } else if (userItem.recipe_limit === -1) {
      setLimitValue('unlimited');
    } else if ([0, 5, 10, 20, 50, 100].includes(userItem.recipe_limit)) {
      // If it's a predefined value, use it
      setLimitValue(userItem.recipe_limit.toString());
    } else {
      // Otherwise, it's a custom value
      setLimitValue('custom');
      setShowCustomInput(true);
      setCustomLimitValue(userItem.recipe_limit.toString());
    }
  };

  const handleSaveLimit = async (userId: string) => {
    try {
      let limitToSave: number | null;

      if (limitValue === 'global') {
        limitToSave = null;
      } else if (limitValue === 'unlimited') {
        limitToSave = -1;
      } else if (limitValue === 'custom') {
        // Use the custom input value
        const parsed = parseInt(customLimitValue, 10);
        if (isNaN(parsed) || parsed < 0) {
          toast.error('Invalid custom limit value. Must be 0 or greater.');
          return;
        }
        limitToSave = parsed;
      } else {
        // Use the dropdown value
        const parsed = parseInt(limitValue, 10);
        if (isNaN(parsed) || parsed < 0) {
          toast.error('Invalid limit value');
          return;
        }
        limitToSave = parsed;
      }

      await apiClient.updateUserRecipeLimit(userId, limitToSave);
      toast.success('Recipe limit updated successfully');
      setEditingLimitUserId(null);
      setShowCustomInput(false);
      setCustomLimitValue('');
      // Refresh data
      fetchData();
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to update limit');
    }
  };

  const handleCancelEdit = () => {
    setEditingLimitUserId(null);
    setLimitValue('');
  };

  const formatRecipeLimit = (limit: number | null | undefined): string => {
    if (limit === null || limit === undefined) return 'Global';
    if (limit === -1) return 'Unlimited';
    if (limit === 0) return 'Blocked';
    return limit.toString();
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const formatRelativeDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - date.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays} days ago`;
    if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`;
    if (diffDays < 365) return `${Math.floor(diffDays / 30)} months ago`;
    return `${Math.floor(diffDays / 365)} years ago`;
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 p-6">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">Admin Portal</h1>
        <div className="animate-pulse space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="bg-white dark:bg-gray-800 rounded-lg shadow p-6 h-32"></div>
            ))}
          </div>
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6 h-96"></div>
        </div>
      </div>
    );
  }

  if (!stats) return null;

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 p-6">
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">Admin Portal</h1>

      {/* Statistics Dashboard */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {/* Total Users */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">Total Users</h3>
            <span className="text-2xl">👥</span>
          </div>
          <p className="text-3xl font-bold text-gray-900 dark:text-white">{stats.total_users}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {stats.recent_registrations} new this week
          </p>
        </div>

        {/* Total Recipes */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">Total Recipes</h3>
            <span className="text-2xl">📝</span>
          </div>
          <p className="text-3xl font-bold text-gray-900 dark:text-white">{stats.total_recipes}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {stats.average_recipes_per_user.toFixed(1)} per user
          </p>
        </div>

        {/* Shopping Items */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">Shopping Items</h3>
            <span className="text-2xl">🛒</span>
          </div>
          <p className="text-3xl font-bold text-gray-900 dark:text-white">{stats.total_shopping_items}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">Across all users</p>
        </div>

        {/* Most Active User */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">Most Active</h3>
            <span className="text-2xl">🏆</span>
          </div>
          <p className="text-xl font-bold text-gray-900 dark:text-white truncate">
            {stats.most_active_user || 'No recipes yet'}
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {stats.most_active_user_count} recipes
          </p>
        </div>
      </div>

      {/* System Info */}
      <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4 mb-8 border border-blue-200 dark:border-blue-800">
        <p className="text-sm text-blue-800 dark:text-blue-300">
          <span className="font-semibold">System Started:</span> {formatDate(stats.first_user_date)} (
          {formatRelativeDate(stats.first_user_date)})
        </p>
      </div>

      {/* User Management Table */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">User Management</h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            Manage all registered users and their data
          </p>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  User
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Email
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Joined
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Recipes
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Recipe Limit
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Shopping
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Last Activity
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
              {users.map((userItem) => (
                <tr
                  key={userItem.id}
                  className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors"
                >
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className="text-sm font-medium text-gray-900 dark:text-white">
                        {userItem.username}
                        {userItem.is_admin && (
                          <span className="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-200">
                            Admin
                          </span>
                        )}
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-300">
                    {userItem.email}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-300">
                    {formatDate(userItem.created_at)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-300">
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
                      {userItem.recipe_count}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-300">
                    {editingLimitUserId === userItem.id ? (
                      <div className="flex items-center gap-2">
                        <select
                          value={limitValue}
                          onChange={(e) => {
                            const newValue = e.target.value;
                            setLimitValue(newValue);
                            if (newValue === 'custom') {
                              setShowCustomInput(true);
                              setCustomLimitValue('');
                            } else {
                              setShowCustomInput(false);
                              setCustomLimitValue('');
                            }
                          }}
                          className="text-sm border border-gray-300 dark:border-gray-600 rounded px-2 py-1 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                        >
                          <option value="global">Global</option>
                          <option value="unlimited">Unlimited</option>
                          <option value="0">Blocked</option>
                          <option value="5">5</option>
                          <option value="10">10</option>
                          <option value="20">20</option>
                          <option value="50">50</option>
                          <option value="100">100</option>
                          <option value="custom">Custom...</option>
                        </select>
                        {showCustomInput && (
                          <input
                            type="number"
                            value={customLimitValue}
                            onChange={(e) => setCustomLimitValue(e.target.value)}
                            min="0"
                            className="text-sm border border-gray-300 dark:border-gray-600 rounded px-2 py-1 bg-white dark:bg-gray-700 text-gray-900 dark:text-white w-20"
                            placeholder="Enter"
                            autoFocus
                          />
                        )}
                        <button
                          onClick={() => handleSaveLimit(userItem.id)}
                          className="text-green-600 hover:text-green-800 dark:text-green-400 dark:hover:text-green-300"
                          title="Save"
                        >
                          ✓
                        </button>
                        <button
                          onClick={handleCancelEdit}
                          className="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                          title="Cancel"
                        >
                          ✗
                        </button>
                      </div>
                    ) : (
                      <div className="flex items-center gap-2">
                        <span
                          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                            userItem.recipe_limit === 0
                              ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                              : userItem.recipe_limit === -1
                              ? 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200'
                              : userItem.recipe_limit === null || userItem.recipe_limit === undefined
                              ? 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
                              : 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200'
                          }`}
                        >
                          {formatRecipeLimit(userItem.recipe_limit)}
                        </span>
                        <button
                          onClick={() => handleEditLimit(userItem)}
                          className="text-primary-600 hover:text-primary-800 dark:text-primary-400 dark:hover:text-primary-300 text-xs"
                          title="Edit limit"
                        >
                          Edit
                        </button>
                      </div>
                    )}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-300">
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                      {userItem.shopping_items}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-300">
                    {userItem.last_recipe_date ? formatDate(userItem.last_recipe_date) : 'Never'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    {userItem.id === user?.id ? (
                      <span
                        className="text-gray-400 dark:text-gray-600 cursor-not-allowed"
                        title="Cannot delete your own account"
                      >
                        Delete
                      </span>
                    ) : (
                      <button
                        onClick={() => handleDeleteClick(userItem)}
                        className="text-primary-600 hover:text-primary-900 dark:text-primary-400 dark:hover:text-primary-300"
                      >
                        Delete
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {users.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500 dark:text-gray-400">No users found</p>
          </div>
        )}
      </div>

      {/* Delete Confirmation Dialog */}
      <ConfirmDialog
        isOpen={showDeleteConfirm}
        title="Delete User?"
        message={`Are you sure you want to delete ${selectedUser?.username}? This will permanently delete their account, all recipes, and shopping list items. This action cannot be undone.`}
        confirmLabel={deleting ? 'Deleting...' : 'Delete User'}
        cancelLabel="Cancel"
        variant="danger"
        onConfirm={handleDeleteConfirm}
        onCancel={() => {
          setShowDeleteConfirm(false);
          setSelectedUser(null);
        }}
      />
    </div>
  );
};
