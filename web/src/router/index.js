import { createRouter, createWebHistory } from 'vue-router'
import { useCookies } from 'vue3-cookies';
import userService from '@/services/userService';
import Placeholder from '../components/Placeholder.vue';

import { store } from '../store';

const { cookies } = useCookies();

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/table/:id',
      name: 'table',
      component: () => import('../views/TableView.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/payment/overview/:id',
      name: 'payment.overview',
      component: Placeholder,
      component: () => import('../views/PaymentOverviewView.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/payment/vipps/:id',
      name: 'payment.vipps',
      component: Placeholder,
      component: () => import('../views/VippsPayment.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/payment/physical/:id',
      name: 'payment.physical',
      component: Placeholder,
      component: () => import('../views/PaymentPhysical.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/payment/result/:id',
      name: 'payment.result',
      component: Placeholder,
      component: () => import('../views/PaymentResult.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/profile/account',
      name: 'profile.account',
      component: () => import('../views/Profile/AccountView.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/profile/order-history',
      name: 'profile.orderHistory',
      component: () => import('../views/Profile/OrderHistoryView.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/profile/waiting-orders',
      name: 'profile.waitingOrders',
      component: () => import('../views/WaitingOrders.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/profile/subscriptions',
      name: 'profile.subscriptions',
      component: () => import('../views/SubscriptionsView.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/profile/change-password',
      name: 'changePassword',
      component: () => import('../views/ChangePasswordView.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/shop',
      name: 'shop',
      component: Placeholder,
      component: () => import('../views/ShopView.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/verify',
      name: 'verify',
      component: () => import('../views/VerifyView.vue'),
      meta: { authOnly: true, noCheck: true }
    },
    {
      path: '/verify/:token',
      name: 'verifyComplete',
      component: () => import('../views/VerifyCompleteView.vue'),
      meta: { authOnly: true, noCheck: true }
    },
    {
      path: '/notifications',
      name: 'notifications',
      component: () => import('../views/NotificationsView.vue'),
      meta: { authOnly: true }
    },
    {
      path: '/admin',
      name: 'admin.dashboard',
      component: () => import('../views/Admin/Dashboard.vue'),
      meta: { authOnly: true, adminOnly: true }
    },

    // Admin - Users
    {
      path: '/admin/users',
      name: 'admin.users',
      component: () => import('../views/Admin/Users.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/users/:id',
      name: 'admin.users.edit',
      component: () => import('../views/Admin/EditUser.vue'),
      meta: { authOnly: true, adminOnly: true }
    },

    // Admin - Tables
    {
      path: '/admin/tables',
      name: 'admin.tables',
      component: () => import('../views/Admin/Tables/Index.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/tables/:id',
      name: 'admin.table',
      component: () => import('../views/Admin/Tables/Edit.vue'),
      meta: { authOnly: true, adminOnly: true }
    },

    // Admin - Reports
    {
      path: '/admin/reports',
      name: 'admin.reports',
      component: () => import('../views/Admin/Reports/Index.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/reports/tables',
      name: 'admin.reports.tables',
      component: () => import('../views/Admin/Reports/Tables.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/reports/games',
      name: 'admin.reports.games',
      component: () => import('../views/Admin/Reports/Games.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/reports/kiosk',
      name: 'admin.reports.kiosk',
      component: () => import('../views/Admin/Reports/Kiosk.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/reports/subscriptions',
      name: 'admin.reports.subscriptions',
      component: () => import('../views/Admin/Reports/Subscriptions.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/reports/games',
      name: 'admin.reports.games',
      component: () => import('../views/Admin/Reports/Games.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/reports/users',
      name: 'admin.reports.users',
      component: () => import('../views/Admin/Reports/Users.vue'),
      meta: { authOnly: true, adminOnly: true }
    },

    // Admin - Roles
    {
      path: '/admin/roles',
      name: 'admin.roles',
      component: () => import('../views/Admin/Roles/Index.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/roles/:id',
      name: 'admin.role',
      component: () => import('../views/Admin/Roles/Edit.vue'),
      meta: { authOnly: true, adminOnly: true }
    },

    // Admin - Products
    {
      path: '/admin/products',
      name: 'admin.products',
      component: () => import('../views/Admin/Products/Index.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/products/:id',
      name: 'admin.product',
      component: () => import('../views/Admin/Products/Edit.vue'),
      meta: { authOnly: true, adminOnly: true }
    },

    // Admin - Pricing
    {
      path: '/admin/pricing',
      name: 'admin.pricing',
      component: () => import('../views/Admin/Pricing/Index.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/pricing/:id',
      name: 'admin.pricing.edit',
      component: () => import('../views/Admin/Pricing/Edit.vue'),
      meta: { authOnly: true, adminOnly: true }
    },

    // Admin - Subscriptions
    {
      path: '/admin/subscriptions',
      name: 'admin.subscriptions',
      component: () => import('../views/Admin/Subscriptions.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/subscriptions/:id',
      name: 'admin.subscription',
      component: () => import('../views/Admin/Subscription.vue'),
      meta: { authOnly: true, adminOnly: true }
    },

    // Admin - Orders
    {
      path: '/admin/orders',
      name: 'admin.orders',
      component: () => import('../views/Admin/Orders.vue'),
      meta: { authOnly: true, adminOnly: true }
    },
    {
      path: '/admin/orders/:id',
      name: 'admin.order',
      component: () => import('../views/Admin/Order.vue'),
      meta: { authOnly: true, adminOnly: true }
    },

    {
      path: '/network-error',
      name: 'networkError',
      component: () => import('../views/NetworkErrorView.vue'),
    },
    {
      path: '/reset-password',
      name: 'resetPassword',
      component: () => import('../views/ResetPasswordView.vue'),
      meta: { guestOnly: true }
    },
    {
      path: '/reset-password-continue/:token',
      name: 'resetPasswordContinue',
      component: () => import('../views/ResetPasswordContinueView.vue'),
      meta: { guestOnly: true }
    },
    {
      path: '/privacy-policy',
      name: 'privacy',
      component: () => import('../views/PrivacyView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { guestOnly: true }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
      meta: { guestOnly: true }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('../views/NotFound.vue'),
    },
  ]
})

router.beforeEach((to, from) => {
  const isAuthenticated = cookies.get('token') !== null;
  
  if (to.meta.noCheck) return true

  if (to.meta.authOnly && !isAuthenticated) return '/login'

  if (to.meta.guestOnly && isAuthenticated) return '/'

  // check is token working
  if (to.meta.authOnly && isAuthenticated && to.name !== 'profile') {
    userService.getProfile()
      .then((response) => {
        store.setUser({
          id: response.data.user.id,
          role_id: response.data.user.role_id,
        });

        if (to.meta.adminOnly && response.data.user.is_admin != true) {
				return '/'
		}
		
        if (response.data.user != null && response.data.user.is_verified === false) {
				return '/verify'
		}

		if (store.user != null && store.user.hasOwnProperty('role_id') && store.user.role_id == 3) {
				console.log(' 3 ')
				if (to.path == '/') {
						return '/shop'
				}

				if (from.path == '/' && to.path != '/shop') {
						return '/shop'
				}
		}

		return
      })
      .catch((err) => {
        console.error(err)
        if (err.response.data.message === "error.user_not_found" || err.response.data.message === "error.auth_required") {
          cookies.remove('token')
          return '/login'
        }
      });

		if (store.user != null && store.user.hasOwnProperty('role_id') && store.user.role_id == 3) {
				if (to.path == '/') {
						return '/shop'
				}

				if (from.path == '/' && to.path != '/shop') {
						return '/shop'
				}
		}
  }
});

export default router
