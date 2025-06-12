   function toggleUserDropdown() {
            const dropdown = document.getElementById('user-dropdown');
            dropdown.classList.toggle('show');
            
            // Close other dropdowns
            document.querySelectorAll('.nav-dropdown-menu').forEach(menu => {
                menu.classList.remove('show');
            });
        }

        function toggleNavDropdown(dropdownId) {
            const dropdown = document.getElementById(dropdownId);
            dropdown.classList.toggle('show');
            
            // Close other dropdowns
            document.getElementById('user-dropdown')?.classList.remove('show');
            document.querySelectorAll('.nav-dropdown-menu').forEach(menu => {
                if (menu.id !== dropdownId) {
                    menu.classList.remove('show');
                }
            });
        }

        // Close dropdowns when clicking outside
        document.addEventListener('click', function(event) {
            const userDropdown = document.getElementById('user-dropdown');
            const userToggle = document.querySelector('.dropdown-toggle');
            
            if (userDropdown && !userDropdown.contains(event.target) && !userToggle?.contains(event.target)) {
                userDropdown.classList.remove('show');
            }
            
            // Close nav dropdowns
            document.querySelectorAll('.nav-dropdown').forEach(navDropdown => {
                const toggle = navDropdown.querySelector('.nav-dropdown-toggle');
                const menu = navDropdown.querySelector('.nav-dropdown-menu');
                
                if (!navDropdown.contains(event.target)) {
                    menu.classList.remove('show');
                }
            });
        });