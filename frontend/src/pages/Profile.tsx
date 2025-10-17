import React, { useState, useEffect } from 'react';
import { useLanguage } from '../contexts/LanguageContext';
import { usePageTitle } from '../hooks/usePageTitle';
import { apiClient } from '../api/client';
import { Input } from '../components/Input';
import type { User } from '../types';

interface UserStats {
  total_recipes: number;
  favorite_recipes: number;
}

export const Profile: React.FC = () => {
  const { t } = useLanguage();
  usePageTitle('Profile');
  const [user, setUser] = useState<User | null>(null);
  const [stats, setStats] = useState<UserStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isEditing, setIsEditing] = useState(false);
  const [username, setUsername] = useState('');
  const [updating, setUpdating] = useState(false);

  useEffect(() => {
    loadProfileData();
  }, []);

  const loadProfileData = async () => {
    try {
      setLoading(true);
      const [profileData, statsData] = await Promise.all([
        apiClient.getProfile(),
        apiClient.getUserStats(),
      ]);
      setUser(profileData);
      setStats(statsData);
      setUsername(profileData.username);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load profile');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!username.trim() || username.length < 3) {
      setError('Username must be at least 3 characters');
      return;
    }

    try {
      setUpdating(true);
      setError('');
      setSuccess('');
      await apiClient.updateProfile({ username });
      setSuccess(t.profile.profileUpdated);
      setIsEditing(false);
      // Reload profile to get updated data
      await loadProfileData();
    } catch (err: any) {
      setError(err.response?.data?.error || t.profile.updateFailed);
    } finally {
      setUpdating(false);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (!user || !stats) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        {error || 'Failed to load profile'}
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">{t.profile.title}</h1>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
          {error}
        </div>
      )}

      {success && (
        <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded mb-6">
          {success}
        </div>
      )}

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-gradient-to-br from-primary-500 to-primary-600 rounded-lg shadow-lg p-6 text-white">
          <div className="text-sm font-medium opacity-90 mb-2">{t.profile.totalRecipes}</div>
          <div className="text-4xl font-bold">{stats.total_recipes}</div>
        </div>
        <div className="bg-gradient-to-br from-yellow-500 to-orange-600 rounded-lg shadow-lg p-6 text-white">
          <div className="text-sm font-medium opacity-90 mb-2">{t.profile.favoriteRecipes}</div>
          <div className="text-4xl font-bold">{stats.favorite_recipes}</div>
        </div>
        <div className="bg-gradient-to-br from-purple-500 to-pink-600 rounded-lg shadow-lg p-6 text-white">
          <div className="text-sm font-medium opacity-90 mb-2">{t.profile.memberSince}</div>
          <div className="text-lg font-bold">{formatDate(user.created_at)}</div>
        </div>
      </div>

      {/* Account Information Card */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-8">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white">{t.profile.accountInfo}</h2>
          {!isEditing && (
            <button
              onClick={() => setIsEditing(true)}
              className="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors"
            >
              {t.profile.editProfile}
            </button>
          )}
        </div>

        {isEditing ? (
          <form onSubmit={handleUpdateProfile} className="space-y-6">
            <div>
              <Input
                label={t.auth.username}
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder={t.auth.usernamePlaceholder}
                required
                minLength={3}
              />
            </div>

            <div>
              <Input
                label={t.auth.email}
                type="email"
                value={user.email}
                disabled
                className="bg-gray-100 dark:bg-gray-900 cursor-not-allowed"
              />
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">Email cannot be changed</p>
            </div>

            <div className="flex gap-4">
              <button
                type="submit"
                disabled={updating}
                className="px-6 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {updating ? `${t.common.loading}` : t.profile.updateProfile}
              </button>
              <button
                type="button"
                onClick={() => {
                  setIsEditing(false);
                  setUsername(user.username);
                  setError('');
                }}
                className="px-6 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200 rounded-lg transition-colors"
              >
                {t.common.cancel}
              </button>
            </div>
          </form>
        ) : (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t.auth.username}
              </label>
              <p className="text-lg text-gray-900 dark:text-white">{user.username}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t.auth.email}
              </label>
              <p className="text-lg text-gray-900 dark:text-white">{user.email}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                {t.profile.memberSince}
              </label>
              <p className="text-lg text-gray-900 dark:text-white">{formatDate(user.created_at)}</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
