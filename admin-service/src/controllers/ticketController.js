// admin-service/controllers/ticketController.js
import { Ticket } from '../models/ticketModel.js';

// GET /admin/tickets
export const getAllTickets = async (req, res) => {
  try {
    const tickets = await Ticket.find().sort({ createdAt: -1 });
    res.json(tickets);
  } catch (error) {
    console.error("Error fetching tickets:", error);
    res.status(500).json({ error: 'Failed to fetch tickets' });
  }
};

// GET /admin/tickets/:id
export const getTicketById = async (req, res) => {
  try {
    const ticket = await Ticket.findById(req.params.id);
    if (!ticket) return res.status(404).json({ error: 'Ticket not found' });
    res.json(ticket);
  } catch (error) {
    console.error("Error fetching ticket:", error);
    res.status(500).json({ error: 'Error fetching ticket' });
  }
};

// POST /admin/tickets/:id/respond
export const respondToTicket = async (req, res) => {
  const { message } = req.body;
  if (!message) return res.status(400).json({ error: 'Message is required' });

  try {
    const ticket = await Ticket.findById(req.params.id);
    if (!ticket) return res.status(404).json({ error: 'Ticket not found' });

    ticket.responses.push({
      adminId: req.user?.sub || 'admin-system',
      message,
    });

    await ticket.save();
    res.json({ message: 'Response added', ticket });
  } catch (error) {
    console.error("Error responding to ticket:", error);
    res.status(500).json({ error: 'Failed to respond to ticket' });
  }
};

// POST /admin/tickets/:id/close
export const closeTicket = async (req, res) => {
  try {
    const ticket = await Ticket.findByIdAndUpdate(
      req.params.id,
      { status: 'closed' },
      { new: true }
    );
    if (!ticket) return res.status(404).json({ error: 'Ticket not found' });

    res.json({ message: 'Ticket closed', ticket });
  } catch (error) {
    console.error("Error closing ticket:", error);
    res.status(500).json({ error: 'Failed to close ticket' });
  }
};
