import { expect, test } from '../../fixtures/user-sessions';
import { AdminTeamPage } from '../../fixtures/admin/team-page';

test.describe('The Admin Team Page', () => {
  test.describe('Unauthenticated user', () => {
    test('redirects to login', async ({ page }) => {
      const adminPage = new AdminTeamPage(page);

      await adminPage.goto('bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa');

      const title = adminPage.page.locator('[data-formtitle="login"]');
      await expect(title).toHaveText('Login');
    });
  });

  test.describe('Guest user', () => {
    test('redirects to landing', async ({ guestPage }) => {
      const adminPage = new AdminTeamPage(guestPage.page);

      await adminPage.goto('bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa');

      const title = adminPage.page.locator('h1');
      await expect(title).toHaveText(
        'Thunderdome is an Agile Planning Poker app with a fun theme',
      );
    });
  });

  test.describe('Non Admin Registered User', () => {
    test('redirects to landing', async ({ registeredPage }) => {
      const adminPage = new AdminTeamPage(registeredPage.page);

      await adminPage.goto('bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa');

      const title = adminPage.page.locator('h1');
      await expect(title).toHaveText(
        'Thunderdome is an Agile Planning Poker app with a fun theme',
      );
    });
  });

  test.describe('Admin User', () => {
    test('loads Team page', async ({ registeredPage, adminPage }) => {
      const ap = new AdminTeamPage(adminPage.page);
      const testTeamName = 'E2E TEST ADMIN TEAM';
      const team = await registeredPage.createTeam(testTeamName);

      await ap.goto(team.id);
      await expect(ap.page.locator('h1')).toContainText(testTeamName);
    });
  });
});
