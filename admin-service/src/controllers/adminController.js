// admin-service/controllers/adminController.js
import pool from '../utils/db.js';

// Get all users
export const getAllUsers = async (req, res) => {
  try {
    const result = await pool.query(
      'SELECT id, email, username, is_admin, created_at FROM users'
    );
    res.json(result.rows);
  } catch (error) {
    console.error('Error fetching users:', error);
    res.status(500).json({ error: 'Failed to retrieve users' });
  }
};

// Edit user
export const editUser = async (req, res) => {
  const { id } = req.params;
  const { email, username, is_admin } = req.body;

  try {
    const query =
      'UPDATE users SET email = $1, username = $2, is_admin = $3 WHERE id = $4 RETURNING id, email, username, is_admin';
    const values = [email, username, is_admin, id];
    const result = await pool.query(query, values); // 🔁 Changed `db` to `pool`

    if (result.rowCount === 0) {
      return res.status(404).json({ error: 'User not found' });
    }

    res.json({ message: 'User updated successfully', user: result.rows[0] });
  } catch (err) {
    console.error('Error updating user:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
};

// Delete user
export const deleteUser = async (req, res) => {
  const { id } = req.params;

  try {
    const result = await pool.query('DELETE FROM users WHERE id = $1', [id]); // 🔁 Changed `db` to `pool`

    if (result.rowCount === 0) {
      return res.status(404).json({ error: 'User not found' });
    }

    res.json({ message: 'User deleted successfully' });
  } catch (err) {
    console.error('Error deleting user:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
};
